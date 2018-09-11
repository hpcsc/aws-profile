package handlers

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
	"testing"
)

func setupHandler(credentialsName string, configName string) GetHandler {
	app := kingpin.New("some-app", "some description")
	getHandler := NewGetHandler(app)

	testCredentialsPath, _ := filepath.Abs("./test_data/" + credentialsName)
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)
	app.Parse([]string { "get", "--credentials-path", testCredentialsPath, "--config-path", testConfigPath  })

	return getHandler
}

func TestReturnErrorIfCredentialsFileNotFound(t *testing.T) {
	getHandler := setupHandler("credentials_not_exists", "get_profile_in_neither_file-config")

	success, output := getHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS credentials file")
}

func TestReturnErrorIfConfigFileNotFound(t *testing.T) {
	getHandler := setupHandler("get_profile_in_neither_file-credentials", "config_not_exists")

	success, output := getHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS config file")
}

func TestConfigProfileHasPriorityOverCredentialsProfile(t *testing.T) {
	getHandler := setupHandler("get_config_priority_over_credentials-credentials", "get_config_priority_over_credentials-config")

	success, output := getHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, output, "assuming profile two")
}

func TestReturnEmptyIfProfileInNeitherConfigNorCredentials(t *testing.T) {
	getHandler := setupHandler("get_profile_in_neither_file-credentials", "get_profile_in_neither_file-config")

	success, output := getHandler.Handle()

	assert.True(t, success)
	assert.Equal(t, "", output)
}

func TestReturnCredentialsProfileIfNotFoundInConfig(t *testing.T) {
	getHandler := setupHandler("get_profile_not_in_config-credentials", "get_profile_not_in_config-config")

	success, output := getHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, output, "two_credentials")
}
