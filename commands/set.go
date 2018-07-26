package commands

import (
	"flag"
	"fmt"
		)

type SetCommand struct {
	Command *flag.FlagSet
	Flags   SetCommandFlags
}

type SetCommandFlags struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
}

func NewSetCommand() SetCommand {
	command := flag.NewFlagSet("set", flag.ExitOnError)

	credentialsFilePath := command.String("credentials-path", "~/.aws/credentials", "Path to AWS Credentials file")
	configFilePath := command.String("config-path", "~/.aws/config", "Path to AWS Config file")

	return SetCommand {
		Command: command,
		Flags:   SetCommandFlags{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
		},
	}
}

func (setCommand SetCommand) Handle(arguments []string) {
	command := setCommand.Command
	command.Parse(arguments)
	if command.Parsed() {
		credentialsPath := ExpandHomeDirectory(*setCommand.Flags.CredentialsFilePath)
		configPath := ExpandHomeDirectory(*setCommand.Flags.ConfigFilePath)

		fmt.Printf("%s %s\n", credentialsPath, configPath)
	}
}
