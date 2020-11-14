package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
)

func stubGlobalArgumentsForGetRegion(configName string) GlobalArguments {
	testConfigPath, _ := filepath.Abs("./test_data/" + configName)

	return GlobalArguments{
		ConfigFilePath: testConfigPath,
	}
}

func setupGetRegionHandler() GetRegionHandler {
	app := kingpin.New("some-app", "some description")
	handler := NewGetRegionHandler(app)

	if _, err := app.Parse([]string{"get-region"}); err != nil {
		fmt.Printf("failed to setup test get region handler: %v\n", err)
		os.Exit(1)
	}

	return handler
}

func TestGetRegionHandler(t *testing.T) {
	t.Run("return error if config file not found", func(t *testing.T) {
		handler := setupGetRegionHandler()
		globalArguments := stubGlobalArgumentsForGetRegion("config_not_exists")

		success, output := handler.Handle(globalArguments)

		require.False(t, success)
		require.Contains(t, output, "Fail to read AWS config file")
	})

	t.Run("return 'no region set' when no region set in default profile", func(t *testing.T) {
		handler := setupGetRegionHandler()
		globalArguments := stubGlobalArgumentsForGetRegion("set-config")

		success, output := handler.Handle(globalArguments)

		require.True(t, success)
		require.Equal(t, "no region set", output)
	})

	t.Run("return 'no region set' when region is set with empty value in default profile", func(t *testing.T) {
		handler := setupGetRegionHandler()
		globalArguments := stubGlobalArgumentsForGetRegion("get-region-empty-region-key-config")

		success, output := handler.Handle(globalArguments)

		require.True(t, success)
		require.Equal(t, "no region set", output)
	})

	t.Run("return region set in default profile", func(t *testing.T) {
		handler := setupGetRegionHandler()
		globalArguments := stubGlobalArgumentsForGetRegion("get-region-config")

		success, output := handler.Handle(globalArguments)

		require.True(t, success)
		require.Equal(t, "us-east-1", output)
	})
}
