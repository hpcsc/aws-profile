package handlers

import (
	"github.com/hpcsc/aws-profile-utils/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
	"path/filepath"
	"strings"
	"testing"
)

func noopWriteToFileMock(_ *ini.File, _ string) {
	// noop
}

func setupSetHandler(credentialsName string, configName string, selectProfileFn utils.SelectProfileFn, writeToFileFn utils.WriteToFileFn) SetHandler {
	app := kingpin.New("some-app", "some description")
	setHandler := NewSetHandler(app, selectProfileFn, writeToFileFn)

	testCredentialsPath, _ := filepath.Abs("./test_data/" + credentialsName)
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)
	app.Parse([]string { "set", "--credentials-path", testCredentialsPath, "--config-path", testConfigPath  })

	return setHandler
}

func TestSetHandlerReturnErrorIfCredentialsFileNotFound(t *testing.T) {
	setHandler := setupSetHandler(
		"credentials_not_exists",
		"get_profile_in_neither_file-config",
		nil,
		nil,
		)

	success, output := setHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS credentials file")
}

func TestSetHandlerReturnErrorIfConfigFileNotFound(t *testing.T) {
	setHandler := setupSetHandler(
		"get_profile_in_neither_file-credentials",
		"config_not_exists",
		nil,
		nil,
		)

	success, output := setHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS config file")
}

func TestSelectProfileIsInvokedWithProfileNamesFromBothConfigs(t *testing.T) {
	called := false

	selectProfileMock := func (profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		assert.ElementsMatch(
			t,
			profiles.GetAllDisplayProfileNames(),
			[]string {
				"credentials_profile_1",
				"credentials_profile_2",
				"assume profile config_profile_1",
				"assume profile config_profile_2",
			},
			)

		called = true
		return []byte("credentials_profile_2"), nil
	}

	setHandler := setupSetHandler(
		"set-credentials",
		"set-config",
		selectProfileMock,
		noopWriteToFileMock,
		)

	success, _ := setHandler.Handle()

	assert.True(t, success)
	if !called {
		t.Errorf("selectProfileFn is not invoked")
	}
}

func TestReturnErrorIfSelectedProfileNotInBothConfigAndCredentials(t *testing.T) {
	selectProfileMock := func (profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("a_random_profile"), nil
	}

	setHandler := setupSetHandler(
		"set-credentials",
		"set-config",
		selectProfileMock,
		noopWriteToFileMock,
	)

	success, message := setHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, message, "not found in either credentials or config file")
}

func TestDefaultProfileInCredentialsIsSetCorrectlyWhenCredentialsProfileSelected(t *testing.T) {
	selectProfileMock := func (profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("credentials_profile_2"), nil
	}

	writeToFileMock := func (file *ini.File, unexpandedFilePath string) {
		// when profile is from credentials file, it should set default in credentials file and clear default in config file
		if strings.Contains(unexpandedFilePath, "-credentials") {
			assert.Equal(t, "4", file.Section("default").Key("aws_access_key_id").Value())
			assert.Equal(t, "4", file.Section("default").Key("aws_secret_access_key").Value())
		} else if strings.Contains(unexpandedFilePath, "-config") {
			assert.Equal(t, "", file.Section("default").Key("role_arn").Value())
			assert.Equal(t, "", file.Section("default").Key("source_profile").Value())
		} else {
			assert.Fail(t, "unexpected call to writeToFile")
		}
	}

	setHandler := setupSetHandler(
		"set-credentials",
		"set-config",
		selectProfileMock,
		writeToFileMock,
	)

	success, message := setHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, message, "credentials from profile [credentials_profile_2]")
}

func TestDefaultProfileInConfigIsSetCorrectlyWhenConfigProfileSelected(t *testing.T) {
	selectProfileMock := func (profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("profile config_profile_2"), nil
	}

	writeToFileMock := func (file *ini.File, unexpandedFilePath string) {
		// when profile is from config file, it should set default in config file and keep default in credentials unchanged
		if strings.Contains(unexpandedFilePath, "-config") {
			assert.Equal(t, "2", file.Section("default").Key("role_arn").Value())
			assert.Equal(t, "2", file.Section("default").Key("source_profile").Value())
		} else {
			assert.Fail(t, "unexpected call to writeToFile")
		}
	}

	setHandler := setupSetHandler(
		"set-credentials",
		"set-config",
		selectProfileMock,
		writeToFileMock,
	)

	success, message := setHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, message, "configs from assumed [profile config_profile_2]")
}