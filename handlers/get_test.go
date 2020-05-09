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

type NullLogger struct {
}

func (l *NullLogger) Debugf(format string, args ...interface{}) {
}

func (l *NullLogger) Infof(format string, args ...interface{}) {
}

func (l *NullLogger) Warnf(format string, args ...interface{}) {
}

func (l *NullLogger) Errorf(format string, args ...interface{}) {
}

func (l *NullLogger) Fatalf(format string, args ...interface{}) {
}

func setupHandler() GetHandler {
	app := kingpin.New("some-app", "some description")
	getHandler := NewGetHandler(app, &NullLogger{}, stubGetAWSCallerIdentity, stubReadCachedCallerIdentity, stubWriteCachedCallerIdentity)

	app.Parse([]string{"get"})

	return getHandler
}

func TestGetHandler(t *testing.T) {
	t.Run("return error if credentials file not found", func(t *testing.T) {
		getHandler := setupHandler()
		globalArguments := stubGlobalArgumentsForGet("credentials_not_exists", "get_profile_in_neither_file-config")

		success, output := getHandler.Handle(globalArguments)

		assert.False(t, success)
		assert.Contains(t, output, "Fail to read AWS credentials file")

	})

	t.Run("return error if config file not found", func(t *testing.T) {
		getHandler := setupHandler()
		globalArguments := stubGlobalArgumentsForGet("get_profile_in_neither_file-credentials", "config_not_exists")

		success, output := getHandler.Handle(globalArguments)

		assert.False(t, success)
		assert.Contains(t, output, "Fail to read AWS config file")

	})

	t.Run("return config profile if both config and credentials default are set", func(t *testing.T) {
		getHandler := setupHandler()
		globalArguments := stubGlobalArgumentsForGet("get_config_priority_over_credentials-credentials", "get_config_priority_over_credentials-config")

		success, output := getHandler.Handle(globalArguments)

		assert.True(t, success)
		assert.Contains(t, output, "profile two")

	})

	t.Run("return empty if profile is not in config or credentials", func(t *testing.T) {
		getHandler := setupHandler()
		globalArguments := stubGlobalArgumentsForGet("get_profile_in_neither_file-credentials", "get_profile_in_neither_file-config")

		success, output := getHandler.Handle(globalArguments)

		assert.True(t, success)
		assert.Equal(t, "", output)

	})

	t.Run("return profile from credentials file if config profile is not set", func(t *testing.T) {
		getHandler := setupHandler()
		globalArguments := stubGlobalArgumentsForGet("get_profile_not_in_config-credentials", "get_profile_not_in_config-config")

		success, output := getHandler.Handle(globalArguments)

		assert.True(t, success)
		assert.Contains(t, output, "two_credentials")

	})

	t.Run("return caller identity result if credentials environment variables are set", func(t *testing.T) {
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

	})

	t.Run("return error code if failing to get caller identity from AWS", func(t *testing.T) {
		var testInputs = []struct {
			awsError       string
			expectedOutput string
		}{
			{
				"An error occurred (ExpiredToken) when calling the GetCallerIdentity operation: The security token included in the request is expired",
				"ExpiredToken",
			},
			{
				"An error occurred (InvalidClientTokenId) when calling the GetCallerIdentity operation: The security token included in the request is invalid",
				"InvalidClientTokenId",
			},
		}

		for _, tt := range testInputs {
			t.Run(tt.expectedOutput, func(t *testing.T) {
				app := kingpin.New("some-app", "some description")
				getHandler := NewGetHandler(app, &NullLogger{}, func() (string, error) {
					return "", errors.New(tt.awsError)
				}, stubReadCachedCallerIdentity, stubWriteCachedCallerIdentity)
				app.Parse([]string{"get"})

				os.Setenv("AWS_ACCESS_KEY_ID", "aws-access-key-id")
				os.Setenv("AWS_SECRET_ACCESS_KEY", "aws-secret-access-key")
				os.Setenv("AWS_SESSION_TOKEN", "aws-session-key")
				globalArguments := stubGlobalArgumentsForGet("get_profile_not_in_config-credentials", "get_profile_not_in_config-config")

				success, output := getHandler.Handle(globalArguments)

				assert.True(t, success)
				assert.Contains(t, output, tt.expectedOutput)

				os.Unsetenv("AWS_ACCESS_KEY_ID")
				os.Unsetenv("AWS_SECRET_ACCESS_KEY")
				os.Unsetenv("AWS_SESSION_TOKEN")
			})
		}
	})

	t.Run("return unknown if failing to parse error code from AWS response", func(t *testing.T) {
		app := kingpin.New("some-app", "some description")
		getHandler := NewGetHandler(app, &NullLogger{}, func() (string, error) {
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
	})
}
