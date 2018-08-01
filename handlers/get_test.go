package handlers

import (
	"testing"
	"gopkg.in/alecthomas/kingpin.v2"
				"github.com/stretchr/testify/assert"
	"path/filepath"
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
	getHandler := setupHandler("credentials_not_exists", "2_config")

	success, output := getHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS credentials file")
}

func TestReturnErrorIfConfigFileNotFound(t *testing.T) {
	getHandler := setupHandler("2_credentials", "config_not_exists")

	success, output := getHandler.Handle()

	assert.False(t, success)
	assert.Contains(t, output, "Fail to read AWS config file")
}

func TestConfigProfileHasPriorityOverCredentialsProfile(t *testing.T) {
	getHandler := setupHandler("1_credentials", "1_config")

	success, output := getHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, output, "assuming profile two")
}

func TestReturnEmptyIfProfileInNeitherConfigNorCredentials(t *testing.T) {
	getHandler := setupHandler("2_credentials", "2_config")

	success, output := getHandler.Handle()

	assert.True(t, success)
	assert.Equal(t, "", output)
}

func TestReturnCredentialsProfileIfNotFoundInConfig(t *testing.T) {
	getHandler := setupHandler("3_credentials", "3_config")

	success, output := getHandler.Handle()

	assert.True(t, success)
	assert.Contains(t, output, "two_credentials")
}
