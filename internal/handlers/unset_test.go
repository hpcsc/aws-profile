package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"testing"
)

func setupUnsetHandler(isWindows bool) UnsetHandler {
	app := kingpin.New("some-app", "some description")
	unsetHandler := NewUnsetHandler(app, isWindows)

	if _, err := app.Parse([]string{"unset"}); err != nil {
		fmt.Printf("failed to setup test set handler: %v\n", err)
		os.Exit(1)
	}

	return unsetHandler
}

func TestUnsetHandler(t *testing.T) {
	t.Run("contains unset command for Linux and MacOS in output", func(t *testing.T) {
		unsetHandler := setupUnsetHandler(false)

		success, output := unsetHandler.Handle(GlobalArguments{})

		require.True(t, success)
		require.Equal(t, output, "unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN AWS_REGION AWS_DEFAULT_REGION")
	})

	t.Run("contains unset command for Windows in output", func(t *testing.T) {
		unsetHandler := setupUnsetHandler(true)

		success, output := unsetHandler.Handle(GlobalArguments{})

		require.True(t, success)
		require.Equal(t, output, "Remove-Item Env:\\AWS_ACCESS_KEY_ID, Env:\\AWS_SECRET_ACCESS_KEY, Env:\\AWS_SESSION_TOKEN, Env:\\AWS_REGION, Env:\\AWS_DEFAULT_REGION")
	})
}
