package handlers

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

type VersionHandler struct {
	SubCommand *kingpin.CmdClause
}

func NewVersionHandler(app *kingpin.Application) VersionHandler {
	subCommand := app.Command("version", "show aws-profile version")
	return VersionHandler{
		SubCommand: subCommand,
	}
}

func (handler VersionHandler) Handle(globalArguments GlobalArguments) (bool, string) {
	fmt.Printf("aws-profile (%s)", version.Current())
	return true, ""
}
