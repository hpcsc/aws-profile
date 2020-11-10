package handlers

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
)

func stubGlobalArgumentsForSetRegion(configName string) GlobalArguments {
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)

	return GlobalArguments{
		ConfigFilePath: testConfigPath,
	}
}

func setupSetRegionHandler(selectRegionFn SelectRegionFn, writeToFileFn WriteToFileFn) SetRegionHandler {
	app := kingpin.New("some-app", "some description")
	setRegionHandler := NewSetRegionHandler(app, selectRegionFn, writeToFileFn)

	app.Parse([]string{"set-region"})

	return setRegionHandler
}

func TestSetRegionHandler(t *testing.T) {
	t.Run("return error if config file not found", func(t *testing.T) {
		setRegionHandler := setupSetRegionHandler(nil, nil)
		globalArguments := stubGlobalArgumentsForSetRegion("config_not_exists")

		success, output := setRegionHandler.Handle(globalArguments)

		assert.False(t, success)
		assert.Contains(t, output, "Fail to read AWS config file")

	})

	t.Run("invoke selectRegion with predefined regions and correct title", func(t *testing.T) {
		called := false

		selectRegionMock := func(regions []string, title string) ([]byte, error) {
			assert.Len(t, regions, 25)
			assert.Equal(t, title, "Select an AWS region")

			called = true
			return []byte("us-west-1"), nil
		}

		setRegionHandler := setupSetRegionHandler(selectRegionMock, noopWriteToFileMock)
		globalArguments := stubGlobalArgumentsForSetRegion("set-config")

		success, _ := setRegionHandler.Handle(globalArguments)

		assert.True(t, success)
		if !called {
			t.Errorf("selectRegionFn is not invoked")
		}

	})

	t.Run("set region of default profile in config file", func(t *testing.T) {
		calledWriteToFile := false

		selectRegionMock := func(regions []string, title string) ([]byte, error) {
			return []byte("ap-southeast-2"), nil
		}

		writeToFileMock := func(file *ini.File, unexpandedFilePath string) error {
			calledWriteToFile = true

			// only the config file should be modified and non-region keys
			// in the default profile should not be modified
			if strings.Contains(unexpandedFilePath, "-config") {
				assert.Equal(t, "1", file.Section("default").Key("role_arn").Value())
				assert.Equal(t, "1", file.Section("default").Key("source_profile").Value())
				assert.Equal(t, "ap-southeast-2", file.Section("default").Key("region").Value())
			} else {
				assert.Fail(t, "unexpected call to writeToFile")
			}

			return nil
		}

		setRegionHandler := setupSetRegionHandler(selectRegionMock, writeToFileMock)
		globalArguments := stubGlobalArgumentsForSetRegion("set-config")

		success, message := setRegionHandler.Handle(globalArguments)

		assert.True(t, success)
		assert.Contains(t, message, "[region ap-southeast-2] -> [default.region]")
		assert.True(t, calledWriteToFile)
	})

	t.Run("return success when user cancels in the middle of selection", func(t *testing.T) {
		calledWriteToFile := false
		selectRegionMock := func(regions []string, title string) ([]byte, error) {
			return nil, errors.New("cancelled by user")
		}
		writeToFileMock := func(file *ini.File, unexpandedFilePath string) error {
			calledWriteToFile = true
			return nil
		}

		setRegionHandler := setupSetRegionHandler(selectRegionMock, writeToFileMock)
		globalArguments := stubGlobalArgumentsForSetRegion("set-config")

		success, message := setRegionHandler.Handle(globalArguments)

		assert.True(t, success)
		assert.Empty(t, message)
		assert.False(t, calledWriteToFile)
	})
}
