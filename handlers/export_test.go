package handlers

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hpcsc/aws-profile/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
	"testing"
)

func stubbedGetAWSCredentials(_ *utils.AWSProfile) (credentials.Value, error) {
	value := credentials.Value{
		AccessKeyID:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "session-token",
		ProviderName:    "stubbed-provider",
	}

	return value, nil
}

func setupExportHandler(configName string, selectProfileFn utils.SelectProfileFn, getAWSCredentialsFn utils.GetAWSCredentialsFn) ExportHandler {
	app := kingpin.New("some-app", "some description")
	exportHandler := NewExportHandler(app, selectProfileFn, getAWSCredentialsFn)

	testConfigPath, _ := filepath.Abs("./test_data/" + configName)
	app.Parse([]string { "export", "--config-path", testConfigPath  })

	return exportHandler
}

func TestExportHandlerReturnErrorIfConfigFileNotFound(t *testing.T) {
	exportHandler := setupExportHandler(
		"config_not_exists",
		nil,
		nil,
		)

	success, output := exportHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS config file")
}

func TestSelectProfileIsInvokedWithProfileNamesFromConfigFileOnly(t *testing.T) {
	called := false

	selectProfileMock := func (profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		assert.ElementsMatch(
			t,
			profiles.GetAllDisplayProfileNames(),
			[]string {
				"assume profile config_profile_1",
				"assume profile config_profile_2",
			},
			)

		called = true
		return []byte("profile config_profile_1"), nil
	}

	exportHandler := setupExportHandler(
		"set-config",
		selectProfileMock,
		stubbedGetAWSCredentials,
	)

	success, _ := exportHandler.Handle()

	assert.True(t, success)
	if !called {
		t.Errorf("selectProfileFn is not invoked")
	}
}

func TestOutputContainsExportInstructionForLinuxAndMacOS(t *testing.T) {
	selectProfileMock := func (profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("profile config_profile_1"), nil
	}

	exportHandler := setupExportHandler(
		"set-config",
		selectProfileMock,
		stubbedGetAWSCredentials,
	)

	success, output := exportHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, output, "export AWS_ACCESS_KEY_ID='access-key-id' AWS_SECRET_ACCESS_KEY='secret-access-key' AWS_SESSION_TOKEN='session-token'")
}

func TestOutputContainsExportInstructionForWindows(t *testing.T) {
	selectProfileMock := func (profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("profile config_profile_1"), nil
	}

	exportHandler := setupExportHandler(
		"set-config",
		selectProfileMock,
		stubbedGetAWSCredentials,
	)

	success, output := exportHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, output, "$env:AWS_ACCESS_KEY_ID = 'access-key-id'; $env:AWS_SECRET_ACCESS_KEY = 'secret-access-key'; $env:AWS_SESSION_TOKEN = 'session-token'")
}
