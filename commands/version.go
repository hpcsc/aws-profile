package commands

import (
	"flag"
	"fmt"
		)

type VersionCommand struct {
	Command *flag.FlagSet
}

func NewVersionCommand() VersionCommand {
	command := flag.NewFlagSet("version", flag.ExitOnError)

	return VersionCommand {
		Command: command,
	}
}

var version = "undefined"
func (versionCommand VersionCommand) Handle(arguments []string) {
	fmt.Printf("aws-profile-utils (v%s)", version)
}
