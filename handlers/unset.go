package handlers

import (
	"github.com/hpcsc/aws-profile/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type UnsetHandler struct {
	SubCommand        *kingpin.CmdClause
	IsWindows         bool
}

func NewUnsetHandler(app *kingpin.Application, isWindows bool) UnsetHandler {
	subCommand := app.Command("unset", `print commands to unset AWS credentials environment variables

To execute the command without printing it to console:

- For Linux/MacOS, execute: "eval $(aws-profile unset)"

- For Windows, execute: "Invoke-Expression (path\to\aws-profile.exe unset)"`)

	return UnsetHandler{
		SubCommand:        subCommand,
		IsWindows:         isWindows,
	}
}

func (handler UnsetHandler) Handle(_ utils.GlobalArguments) (bool, string) {
	output := formatUnsetCommandByPlatform(handler.IsWindows)
	return true, output
}

func formatUnsetCommandByPlatform(isWindows bool) string {
	if isWindows {
		return "Remove-Item Env:\\AWS_ACCESS_KEY_ID, Env:\\AWS_SECRET_ACCESS_KEY, Env:\\AWS_SESSION_TOKEN"
	}

	return "unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN"
}
