package commands

import (
	"gopkg.in/ini.v1"
	"fmt"
	"os"
	"strings"
	"flag"
)

type GetCommand struct {
	Command *flag.FlagSet
	Flags GetCommandFlags
}

type GetCommandFlags struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
}

func NewGetCommand() GetCommand {
	command := flag.NewFlagSet("get", flag.ExitOnError)

	credentialsFilePath := command.String("credentials-path", "~/.aws/credentials", "Path to AWS Credentials file")
	configFilePath := command.String("config-path", "~/.aws/config", "Path to AWS Config file")

	return GetCommand {
		Command: command,
		Flags:   GetCommandFlags{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
		},
	}
}

func (getCommand GetCommand) Handle(arguments []string) {
	command := getCommand.Command
	command.Parse(arguments)
	if command.Parsed() {
		credentialsPath := ExpandHomeDirectory(*getCommand.Flags.CredentialsFilePath)
		configPath := ExpandHomeDirectory(*getCommand.Flags.ConfigFilePath)

		configFile, err := ini.Load(configPath)
		if err != nil {
			fmt.Printf("Fail to read AWS config file: %v", err)
			os.Exit(1)
		}

		configDefaultSection, err := configFile.GetSection("default")
		if err == nil &&
			configDefaultSection.HasKey("role_arn") &&
			configDefaultSection.HasKey("source_profile") {

			defaultRoleArn := configDefaultSection.Key("role_arn").Value()
			defaultSourceProfile := configDefaultSection.Key("source_profile").Value()

			for _, section := range configFile.Sections() {
				if strings.Compare(section.Name(), "default") != 0 &&
					section.Haskey("role_arn") &&
					section.HasKey("source_profile") &&
					strings.Compare(section.Key("role_arn").Value(), defaultRoleArn) == 0 &&
					strings.Compare(section.Key("source_profile").Value(), defaultSourceProfile) == 0 {
					fmt.Printf("%s", section.Name())
					os.Exit(0)
				}
			}
		}

		credentialsFile, err := ini.Load(credentialsPath)
		if err != nil {
			fmt.Printf("Fail to read AWS credentials file: %v", err)
			os.Exit(1)
		}

		credentialsDefaultSection, err := credentialsFile.GetSection("default")
		if err == nil &&
			credentialsDefaultSection.HasKey("aws_access_key_id") {

			defaultAWSAccessKeyId := credentialsDefaultSection.Key("aws_access_key_id").Value()

			for _, section := range credentialsFile.Sections() {
				if strings.Compare(section.Name(), "default") != 0 &&
					section.HasKey("aws_access_key_id") &&
					strings.Compare(section.Key("aws_access_key_id").Value(), defaultAWSAccessKeyId) == 0 {
					fmt.Printf("%s", section.Name())
					os.Exit(0)
				}
			}
		}
	}
}
