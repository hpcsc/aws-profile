package handlers

import (
	"errors"
	"fmt"
	"github.com/hpcsc/aws-profile/internal/config"
	"github.com/hpcsc/aws-profile/internal/utils"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func noopWriteToFileMock(_ *ini.File, _ string) error {
	// noop
	return nil
}

func stubGlobalArgumentsForSet(credentialsName string, configName string) GlobalArguments {
	testCredentialsPath, _ := filepath.Abs("./test_data/" + credentialsName)
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)

	return GlobalArguments{
		CredentialsFilePath: testCredentialsPath,
		ConfigFilePath:      testConfigPath,
	}
}

func setupSetHandler(selectProfileFn SelectProfileFn, writeToFileFn WriteToFileFn) SetHandler {
	app := kingpin.New("some-app", "some description")
	setHandler := NewSetHandler(app, selectProfileFn, writeToFileFn)

	if _, err := app.Parse([]string{"set"}); err != nil {
		fmt.Printf("failed to setup test set handler: %v\n", err)
		os.Exit(1)
	}

	return setHandler
}

func TestSetHandler(t *testing.T) {
	t.Run("return error if credentials file not found", func(t *testing.T) {
		setHandler := setupSetHandler(nil, nil)
		globalArguments := stubGlobalArgumentsForSet("credentials_not_exists", "get_profile_in_neither_file-config")

		success, output := setHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, output, "Fail to read AWS credentials file")

	})

	t.Run("return error if config file not found", func(t *testing.T) {
		setHandler := setupSetHandler(nil, nil)
		globalArguments := stubGlobalArgumentsForSet("get_profile_in_neither_file-credentials", "config_not_exists")

		success, output := setHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, output, "Fail to read AWS config file")

	})

	t.Run("invoke selectProfile with profile names from both credentials and config files", func(t *testing.T) {
		called := false

		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			require.ElementsMatch(
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

		require.True(t, success)
		if !called {
			t.Errorf("selectProfileFn is not invoked")
		}

	})

	t.Run("return error if selected profile not found in both config and credentials file", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("a_random_profile"), nil
		}

		setHandler := setupSetHandler(selectProfileMock, noopWriteToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, message, "not found in either credentials or config file")
	})

	t.Run("return success when user cancels in the middle of selection", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return nil, utils.NewCancelledError()
		}

		setHandler := setupSetHandler(selectProfileMock, noopWriteToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		require.True(t, success)
		require.Empty(t, message)
	})

	t.Run("set default profile in credentials file when profile is in credentials file", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("credentials_profile_2"), nil
		}

		writeToFileMock := func(file *ini.File, unexpandedFilePath string) error {
			// when profile is from credentials file, it should set default in credentials file and clear default in config file
			if strings.Contains(unexpandedFilePath, "-credentials") {
				require.Equal(t, "4", file.Section("default").Key("aws_access_key_id").Value())
				require.Equal(t, "4", file.Section("default").Key("aws_secret_access_key").Value())
			} else if strings.Contains(unexpandedFilePath, "-config") {
				require.Equal(t, "", file.Section("default").Key("role_arn").Value())
				require.Equal(t, "", file.Section("default").Key("source_profile").Value())
			} else {
				require.Fail(t, "unexpected call to writeToFile")
			}

			return nil
		}

		setHandler := setupSetHandler(selectProfileMock, writeToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		require.True(t, success)
		require.Contains(t, message, "[credentials_profile_2] -> [default]")
	})

	t.Run("return error when profile is in credentials file and failed to write updated credentials file", func(t *testing.T) {
		selectProfileStub := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("credentials_profile_2"), nil
		}

		writeToFileStub := func(file *ini.File, unexpandedFilePath string) error {
			// when profile is from credentials file, it should set default in credentials file and clear default in config file
			if strings.Contains(unexpandedFilePath, "-credentials") {
				return errors.New("some error")
			}

			return nil
		}

		setHandler := setupSetHandler(selectProfileStub, writeToFileStub)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, message, "some error")
	})

	t.Run("return error when profile is in credentials file and failed to write updated config file", func(t *testing.T) {
		selectProfileStub := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("credentials_profile_2"), nil
		}

		writeToFileStub := func(file *ini.File, unexpandedFilePath string) error {
			// when profile is from credentials file, it should set default in credentials file and clear default in config file
			if strings.Contains(unexpandedFilePath, "-config") {
				return errors.New("some error")
			}

			return nil
		}

		setHandler := setupSetHandler(selectProfileStub, writeToFileStub)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, message, "some error")
	})

	t.Run("set default profile in config file when profile is in config file", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_2"), nil
		}

		writeToFileMock := func(file *ini.File, unexpandedFilePath string) error {
			// when profile is from config file, it should set default in config file and keep default in credentials unchanged
			if strings.Contains(unexpandedFilePath, "-config") {
				require.Equal(t, "2", file.Section("default").Key("role_arn").Value())
				require.Equal(t, "2", file.Section("default").Key("source_profile").Value())
			} else {
				require.Fail(t, "unexpected call to writeToFile")
			}

			return nil
		}

		setHandler := setupSetHandler(selectProfileMock, writeToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		require.True(t, success)
		require.Contains(t, message, "[profile config_profile_2] -> [default]")
	})

	t.Run("return error when profile is in config file and failed to write updated config file", func(t *testing.T) {
		selectProfileMock := func(profiles config.Profiles, pattern string) ([]byte, error) {
			return []byte("profile config_profile_2"), nil
		}

		writeToFileMock := func(file *ini.File, unexpandedFilePath string) error {
			if strings.Contains(unexpandedFilePath, "-config") {
				return errors.New("some error")
			}

			return nil
		}

		setHandler := setupSetHandler(selectProfileMock, writeToFileMock)
		globalArguments := stubGlobalArgumentsForSet("set-credentials", "set-config")

		success, message := setHandler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, message, "some error")
	})
}
