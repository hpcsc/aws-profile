package handlers

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hpcsc/aws-profile/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
	"testing"
)

func stubGetAWSCredentials(_ *utils.AWSProfile) (credentials.Value, error) {
	value := credentials.Value{
		AccessKeyID:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "session-token",
		ProviderName:    "stubbed-provider",
	}

	return value, nil
}

func stubGlobalArgumentsForExport(configName string) utils.GlobalArguments {
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)

	return utils.GlobalArguments{
		ConfigFilePath: &testConfigPath,
	}
}

func setupExportHandler(isWindows bool, selectProfileFn utils.SelectProfileFn, getAWSCredentialsFn utils.GetAWSCredentialsFn) ExportHandler {
	app := kingpin.New("some-app", "some description")
	exportHandler := NewExportHandler(app, isWindows, selectProfileFn, getAWSCredentialsFn)

	app.Parse([]string{"export"})

	return exportHandler
}

func TestExportHandler_ReturnErrorIfConfigFileNotFound(t *testing.T) {
	exportHandler := setupExportHandler(
		false,
		nil,
		nil,
	)
	globalArguments := stubGlobalArgumentsForExport("config_not_exists")

	success, output := exportHandler.Handle(globalArguments)

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS config file")
}

func TestExportHandler_SelectProfileIsInvokedWithProfileNamesFromConfigFileOnly(t *testing.T) {
	called := false

	selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		assert.ElementsMatch(
			t,
			profiles.GetAllDisplayProfileNames(),
			[]string{
				"assume profile config_profile_1",
				"assume profile config_profile_2",
			},
		)

		called = true
		return []byte("profile config_profile_1"), nil
	}

	exportHandler := setupExportHandler(
		false,
		selectProfileMock,
		stubGetAWSCredentials,
	)
	globalArguments := stubGlobalArgumentsForExport("set-config")

	success, _ := exportHandler.Handle(globalArguments)

	assert.True(t, success)
	if !called {
		t.Errorf("selectProfileFn is not invoked")
	}
}

func TestExportHandler_OutputContainsExportInstructionForLinuxAndMacOS(t *testing.T) {
	selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("profile config_profile_1"), nil
	}

	exportHandler := setupExportHandler(
		false,
		selectProfileMock,
		stubGetAWSCredentials,
	)
	globalArguments := stubGlobalArgumentsForExport("set-config")

	success, output := exportHandler.Handle(globalArguments)

	assert.True(t, success)
	assert.Equal(t, output, "export AWS_ACCESS_KEY_ID='access-key-id' AWS_SECRET_ACCESS_KEY='secret-access-key' AWS_SESSION_TOKEN='session-token'")
}

func TestExportHandler_OutputContainsExportInstructionForWindows(t *testing.T) {
	selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("profile config_profile_1"), nil
	}

	exportHandler := setupExportHandler(
		true,
		selectProfileMock,
		stubGetAWSCredentials,
	)
	globalArguments := stubGlobalArgumentsForExport("set-config")

	success, output := exportHandler.Handle(globalArguments)

	assert.True(t, success)
	assert.Equal(t, output, "$env:AWS_ACCESS_KEY_ID = 'access-key-id'; $env:AWS_SECRET_ACCESS_KEY = 'secret-access-key'; $env:AWS_SESSION_TOKEN = 'session-token'")
}
