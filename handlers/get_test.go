package handlers

import (
	"errors"
	"github.com/hpcsc/aws-profile/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"testing"
)

func stubGlobalArgumentsForGet(credentialsName string, configName string) utils.GlobalArguments {
	testCredentialsPath, _ := filepath.Abs("./test_data/" + credentialsName)
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)

	return utils.GlobalArguments{
		CredentialsFilePath: &testCredentialsPath,
		ConfigFilePath:      &testConfigPath,
	}
}

func stubGetAWSCallerIdentity() (string, error) {
	return "caller-identity-profile", nil
}

func stubReadCachedCallerIdentity() (string, error) {
	return "", nil
}

func stubWriteCachedCallerIdentity(_ string) error {
	return nil
}

func setupHandler() GetHandler {
	app := kingpin.New("some-app", "some description")
	getHandler := NewGetHandler(app, stubGetAWSCallerIdentity, stubReadCachedCallerIdentity, stubWriteCachedCallerIdentity)

	app.Parse([]string{"get"})

	return getHandler
}

func TestGetHandler_ReturnErrorIfCredentialsFileNotFound(t *testing.T) {
	getHandler := setupHandler()
	globalArguments := stubGlobalArgumentsForGet("credentials_not_exists", "get_profile_in_neither_file-config")

	success, output := getHandler.Handle(globalArguments)

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS credentials file")
}

func TestGetHandler_ReturnErrorIfConfigFileNotFound(t *testing.T) {
	getHandler := setupHandler()
	globalArguments := stubGlobalArgumentsForGet("get_profile_in_neither_file-credentials", "config_not_exists")

	success, output := getHandler.Handle(globalArguments)

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS config file")
}

func TestGetHandler_ConfigProfileHasPriorityOverCredentialsProfile(t *testing.T) {
	getHandler := setupHandler()
	globalArguments := stubGlobalArgumentsForGet("get_config_priority_over_credentials-credentials", "get_config_priority_over_credentials-config")

	success, output := getHandler.Handle(globalArguments)

	assert.True(t, success)
	assert.Contains(t, output, "profile two")
}

func TestGetHandler_ReturnEmptyIfProfileInNeitherConfigNorCredentials(t *testing.T) {
	getHandler := setupHandler()
	globalArguments := stubGlobalArgumentsForGet("get_profile_in_neither_file-credentials", "get_profile_in_neither_file-config")

	success, output := getHandler.Handle(globalArguments)

	assert.True(t, success)
	assert.Equal(t, "", output)
}

func TestGetHandler_ReturnCredentialsProfileIfNotFoundInConfig(t *testing.T) {
	getHandler := setupHandler()
	globalArguments := stubGlobalArgumentsForGet("get_profile_not_in_config-credentials", "get_profile_not_in_config-config")

	success, output := getHandler.Handle(globalArguments)

	assert.True(t, success)
	assert.Contains(t, output, "two_credentials")
}

func TestGetHandler_ReturnCallerIdentityResultIfCredentialsEnvironmentVariablesAreSet(t *testing.T) {
	getHandler := setupHandler()
	os.Setenv("AWS_ACCESS_KEY_ID", "aws-access-key-id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "aws-secret-access-key")
	os.Setenv("AWS_SESSION_TOKEN", "aws-session-key")
	globalArguments := stubGlobalArgumentsForGet("get_profile_not_in_config-credentials", "get_profile_not_in_config-config")

	success, output := getHandler.Handle(globalArguments)

	assert.True(t, success)
	assert.Contains(t, output, "caller-identity-profile")

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
}

func TestGetHandler_ReturnUnknownIfFailToGetCallerIdentityFromAWS(t *testing.T) {
	app := kingpin.New("some-app", "some description")
	getHandler := NewGetHandler(app, func() (string, error) {
		return "", errors.New("some error from aws")
	}, stubReadCachedCallerIdentity, stubWriteCachedCallerIdentity)
	app.Parse([]string{"get"})


	os.Setenv("AWS_ACCESS_KEY_ID", "aws-access-key-id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "aws-secret-access-key")
	os.Setenv("AWS_SESSION_TOKEN", "aws-session-key")
	globalArguments := stubGlobalArgumentsForGet("get_profile_not_in_config-credentials", "get_profile_not_in_config-config")

	success, output := getHandler.Handle(globalArguments)

	assert.True(t, success)
	assert.Contains(t, output, "unknown")

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
}
