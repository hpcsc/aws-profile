package handlers

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hpcsc/aws-profile/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
	"testing"
	"time"
)

func stubAWSCredentials() credentials.Value {
	return credentials.Value {
		AccessKeyID:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "session-token",
		ProviderName:    "stubbed-provider",
	}
}

func stubGetAWSCredentials(_ *utils.AWSProfile, _ time.Duration) (credentials.Value, error) {
	return stubAWSCredentials(), nil
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

func TestExportHandler_ReturnErrorIfDurationIsInvalid(t *testing.T) {
	app := kingpin.New("some-app", "some description")
	exportHandler := NewExportHandler(app, false, nil, nil)

	app.Parse([]string{"export", "-d", "5"})
	globalArguments := stubGlobalArgumentsForExport("set-config")

	success, err := exportHandler.Handle(globalArguments)

	assert.False(t, success)
	assert.Contains(t, err, "missing unit in duration")
}

func TestExportHandler_ReturnErrorIfDurationIsLowerThanMinimumAllowed(t *testing.T) {
	app := kingpin.New("some-app", "some description")
	exportHandler := NewExportHandler(app, false, nil, nil)

	app.Parse([]string{"export", "-d", "5m"})
	globalArguments := stubGlobalArgumentsForExport("set-config")

	success, err := exportHandler.Handle(globalArguments)

	assert.False(t, success)
	assert.Contains(t, err, "Minimum duration is 15 minutes")
}

func TestExportHandler_CallGetAWSCredentialsWithDefaultValueWhenNoDurationGiven(t *testing.T) {
	selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("profile config_profile_1"), nil
	}

	called := false

	getAWSCredentialsMock := func(_ *utils.AWSProfile, duration time.Duration) (credentials.Value, error) {
		assert.Equal(t, float64(15), duration.Minutes())
		called = true
		return stubAWSCredentials(), nil
	}

	exportHandler := setupExportHandler(
		false,
		selectProfileMock,
		getAWSCredentialsMock,
	)
	globalArguments := stubGlobalArgumentsForExport("set-config")

	exportHandler.Handle(globalArguments)

	assert.True(t, called)
}

func TestExportHandler_CallGetAWSCredentialsWithValueGiven(t *testing.T) {
	selectProfileMock := func(profiles utils.AWSProfiles, pattern string) ([]byte, error) {
		return []byte("profile config_profile_1"), nil
	}

	called := false
	mockDurationValue := "20m"

	getAWSCredentialsMock := func(_ *utils.AWSProfile, duration time.Duration) (credentials.Value, error) {
		assert.Equal(t, float64(20), duration.Minutes())
		called = true
		return stubAWSCredentials(), nil
	}

	app := kingpin.New("some-app", "some description")
	exportHandler := NewExportHandler(app, false, selectProfileMock, getAWSCredentialsMock)

	app.Parse([]string{"export", "-d", mockDurationValue})
	globalArguments := stubGlobalArgumentsForExport("set-config")

	exportHandler.Handle(globalArguments)

	assert.True(t, called)
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
