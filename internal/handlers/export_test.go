package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hpcsc/aws-profile/internal/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func stubAWSCredentials() credentials.Value {
	return credentials.Value{
		AccessKeyID:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "session-token",
		ProviderName:    "stubbed-provider",
	}
}

func stubGetAWSCredentials(_ *config.Profile, _ time.Duration) (credentials.Value, error) {
	return stubAWSCredentials(), nil
}

func stubGlobalArgumentsForExport(configName string) GlobalArguments {
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)

	return GlobalArguments{
		ConfigFilePath: testConfigPath,
	}
}

func setupExportHandler(isWindows bool, selectProfileFn SelectProfileFn, getAWSCredentialsFn GetAWSCredentialsFn) ExportHandler {
	app := kingpin.New("some-app", "some description")
	exportHandler := NewExportHandler(app, isWindows, selectProfileFn, getAWSCredentialsFn)

	if _, err := app.Parse([]string{"export"}); err != nil {
		fmt.Printf("failed to setup test export handler: %v\n", err)
		os.Exit(1)
	}

	return exportHandler
}

func TestExportHandler(t *testing.T) {
	t.Run("return error if config file is not found", func(t *testing.T) {
		exportHandler := setupExportHandler(
			false,
			nil,
			nil,
		)
		globalArguments := stubGlobalArgumentsForExport("config_not_exists")

		success, output := exportHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, output, "Fail to read AWS config file")
	})

	t.Run("invoke SelectProfile with profile names from config file only", func(t *testing.T) {
		called := false

		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			require.ElementsMatch(
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

		require.True(t, success)
		if !called {
			t.Errorf("selectProfileFn is not invoked")
		}
	})

	t.Run("return error if duration is invalid", func(t *testing.T) {
		app := kingpin.New("some-app", "some description")
		exportHandler := NewExportHandler(app, false, nil, nil)

		if _, err := app.Parse([]string{"export", "-d", "5"}); err != nil {
			t.Fatalf("failed to setup test export handler: %v\n", err)
		}

		globalArguments := stubGlobalArgumentsForExport("set-config")

		success, err := exportHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, err, "missing unit in duration")
	})

	t.Run("return error if duration is lower than minimum duration allowed", func(t *testing.T) {
		app := kingpin.New("some-app", "some description")
		exportHandler := NewExportHandler(app, false, nil, nil)

		if _, err := app.Parse([]string{"export", "-d", "5m"}); err != nil {
			t.Fatalf("failed to setup test export handler: %v\n", err)
		}

		globalArguments := stubGlobalArgumentsForExport("set-config")

		success, err := exportHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, err, "Minimum duration is 15 minutes")
	})

	t.Run("call GetAWSCredentials with default value when no duration given", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_1"), nil
		}

		called := false

		getAWSCredentialsMock := func(_ *config.Profile, duration time.Duration) (credentials.Value, error) {
			require.Equal(t, float64(15), duration.Minutes())
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

		require.True(t, called)
	})

	t.Run("call GetAWSCredentials with given value", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_1"), nil
		}

		called := false
		mockDurationValue := "20m"

		getAWSCredentialsMock := func(_ *config.Profile, duration time.Duration) (credentials.Value, error) {
			require.Equal(t, float64(20), duration.Minutes())
			called = true
			return stubAWSCredentials(), nil
		}

		app := kingpin.New("some-app", "some description")
		exportHandler := NewExportHandler(app, false, selectProfileMock, getAWSCredentialsMock)

		if _, err := app.Parse([]string{"export", "-d", mockDurationValue}); err != nil {
			t.Fatalf("failed to setup test export handler: %v\n", err)
		}

		globalArguments := stubGlobalArgumentsForExport("set-config")

		exportHandler.Handle(globalArguments)

		require.True(t, called)
	})

	t.Run("contains export command for Linux and MacOS in output", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_1"), nil
		}

		exportHandler := setupExportHandler(
			false,
			selectProfileMock,
			stubGetAWSCredentials,
		)
		globalArguments := stubGlobalArgumentsForExport("set-config")

		success, output := exportHandler.Handle(globalArguments)

		require.True(t, success)
		require.Equal(t, output, "export AWS_ACCESS_KEY_ID='access-key-id' AWS_SECRET_ACCESS_KEY='secret-access-key' AWS_SESSION_TOKEN='session-token'")
	})

	t.Run("contains export region for Linux and MacOS in output", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_2"), nil
		}

		exportHandler := setupExportHandler(
			false,
			selectProfileMock,
			stubGetAWSCredentials,
		)
		globalArguments := stubGlobalArgumentsForExport("set-config")

		success, output := exportHandler.Handle(globalArguments)

		require.True(t, success)
		require.Equal(t, output, "export AWS_ACCESS_KEY_ID='access-key-id' AWS_SECRET_ACCESS_KEY='secret-access-key' AWS_SESSION_TOKEN='session-token' AWS_REGION='us-west-2' AWS_DEFAULT_REGION='us-west-2'")
	})

	t.Run("contains export command for Windows in output", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_1"), nil
		}

		exportHandler := setupExportHandler(
			true,
			selectProfileMock,
			stubGetAWSCredentials,
		)
		globalArguments := stubGlobalArgumentsForExport("set-config")

		success, output := exportHandler.Handle(globalArguments)

		require.True(t, success)
		require.Equal(t, output, "$env:AWS_ACCESS_KEY_ID = 'access-key-id'; $env:AWS_SECRET_ACCESS_KEY = 'secret-access-key'; $env:AWS_SESSION_TOKEN = 'session-token'")
	})

	t.Run("contains export region for Windows in output", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_2"), nil
		}

		exportHandler := setupExportHandler(
			true,
			selectProfileMock,
			stubGetAWSCredentials,
		)
		globalArguments := stubGlobalArgumentsForExport("set-config")

		success, output := exportHandler.Handle(globalArguments)

		require.True(t, success)
		require.Equal(t, output, "$env:AWS_ACCESS_KEY_ID = 'access-key-id'; $env:AWS_SECRET_ACCESS_KEY = 'secret-access-key'; $env:AWS_SESSION_TOKEN = 'session-token'; $env:AWS_REGION = 'us-west-2'; $env:AWS_DEFAULT_REGION = 'us-west-2'")
	})
}
