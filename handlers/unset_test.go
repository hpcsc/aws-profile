package handlers

import (
	"github.com/hpcsc/aws-profile/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"testing"
)

func setupUnsetHandler(isWindows bool) UnsetHandler {
	app := kingpin.New("some-app", "some description")
	unsetHandler := NewUnsetHandler(app, isWindows)

	app.Parse([]string{"unset"})

	return unsetHandler
}

func TestUnsetHandler_OutputContainsUnsetInstructionForLinuxAndMacOS(t *testing.T) {
	unsetHandler := setupUnsetHandler(false)

	success, output := unsetHandler.Handle(utils.GlobalArguments{
		CredentialsFilePath: nil,
		ConfigFilePath:      nil,
	})

	assert.True(t, success)
	assert.Equal(t, output, "unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN")
}

func TestUnsetHandler_OutputContainsUnsetInstructionForWindows(t *testing.T) {
	unsetHandler := setupUnsetHandler(true)

	success, output := unsetHandler.Handle(utils.GlobalArguments{
		CredentialsFilePath: nil,
		ConfigFilePath:      nil,
	})

	assert.True(t, success)
	assert.Equal(t, output, "Remove-Item Env:\\AWS_ACCESS_KEY_ID, Env:\\AWS_SECRET_ACCESS_KEY, Env:\\AWS_SESSION_TOKEN")
}
