package handlers

import (
	"github.com/hpcsc/aws-profile/utils"
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

func stubGlobalArgumentsForSet(credentialsName string, configName string) utils.GlobalArguments {
	testCredentialsPath, _ := filepath.Abs("./test_data/" + credentialsName)
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)

	return utils.GlobalArguments{
		CredentialsFilePath: &testCredentialsPath,
		ConfigFilePath:      &testConfigPath,
	}
}

func setupSetHandler(selectProfileFn utils.SelectProfileFn, writeToFileFn utils.WriteToFileFn) SetHandler {
	app := kingpin.New("some-app", "some description")
	setHandler := NewSetHandler(app, selectProfileFn, writeToFileFn)

	app.Parse([]string{"set"})

	return setHandler
}

func TestSetHandler(t *testing.T) {
	t.Run("return error if credentials file not found", func(t *testing.T) {
		setHandler := setupSetHandler(nil, nil)
		globalArguments := stubGlobalArgumentsForSet("credentials_not_exists", "get_profile_in_neither_file-config")

		success, output := setHandler.Handle(globalArguments)

		assert.False(t, success)
		assert.Contains(t, output, "Fail to read AWS credentials file")

	})

	t.Run("return error if config file not found", func(t *testing.T) {
		setHandler := setupSetHandler(nil, nil)
		globalArguments := stubGlobalArgumentsForSet("get_profile_in_neither_file-credentials", "config_not_exists")

		success, output := setHandler.Handle(globalArguments)

		assert.False(t, success)
		assert.Contains(t, output, "Fail to read AWS config file")

	})

	t.Run("invoke selectProfile with profile names from both credentials and config files", func(t *testing.T) {
		called := false

		selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
			assert.ElementsMatch(
				t,
				profiles.GetAllDisplayProfileNames(),
				[]string{
					"credentials_profile_1",
					"credentials_profile_2",
					"assume profile config_profile_1",
					"assume profile config_profile_2",
				},
			)

			called = true
			return []byte("credentials_profile_2"), nil
		}

		setHandler := setupSetHandler(selectProfileMock, noopWriteToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, _ := setHandler.Handle(globalArguments)

		assert.True(t, success)
		if !called {
			t.Errorf("selectProfileFn is not invoked")
		}

	})

	t.Run("return error if selected profile not found in both config and credentials file", func(t *testing.T) {
		selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
			return []byte("a_random_profile"), nil
		}

		setHandler := setupSetHandler(selectProfileMock, noopWriteToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		assert.False(t, success)
		assert.Contains(t, message, "not found in either credentials or config file")

	})

	t.Run("set default profile in credentials file", func(t *testing.T) {
		selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
			return []byte("credentials_profile_2"), nil
		}

		writeToFileMock := func(file *ini.File, unexpandedFilePath string) {
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

		setHandler := setupSetHandler(selectProfileMock, writeToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		assert.True(t, success)
		assert.Contains(t, message, "[credentials_profile_2] -> [default]")

	})

	t.Run("set default profile in config file", func(t *testing.T) {
		selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_2"), nil
		}

		writeToFileMock := func(file *ini.File, unexpandedFilePath string) {
			// when profile is from config file, it should set default in config file and keep default in credentials unchanged
			if strings.Contains(unexpandedFilePath, "-config") {
				assert.Equal(t, "2", file.Section("default").Key("role_arn").Value())
				assert.Equal(t, "2", file.Section("default").Key("source_profile").Value())
			} else {
				assert.Fail(t, "unexpected call to writeToFile")
			}
		}

		setHandler := setupSetHandler(selectProfileMock, writeToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		assert.True(t, success)
		assert.Contains(t, message, "[profile config_profile_2] -> [default]")

	})
}
