package commands

import (
	"flag"
	"fmt"
	"os/exec"
	"os"
	"gopkg.in/ini.v1"
	"strings"
)

type SetCommand struct {
	Command *flag.FlagSet
	Flags   SetCommandFlags
}

type SetCommandFlags struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
	Pattern *string
}

func NewSetCommand() SetCommand {
	command := flag.NewFlagSet("set", flag.ExitOnError)

	credentialsFilePath := command.String("credentials-path", "~/.aws/credentials", "Path to AWS Credentials file")
	configFilePath := command.String("config-path", "~/.aws/config", "Path to AWS Config file")
	pattern := command.String("pattern", "", "Start the fzf finder with the given query")

	return SetCommand {
		Command: command,
		Flags:   SetCommandFlags{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
			Pattern: pattern,
		},
	}
}

func getAWSProfiles(setCommand SetCommand) []string {
	var profiles []string

	credentialsPath := ExpandHomeDirectory(*setCommand.Flags.CredentialsFilePath)
	credentialsFile, err := ini.Load(credentialsPath)
	if err != nil {
		fmt.Printf("Fail to read AWS credentials file: %v", err)
		os.Exit(1)
	}

	for _, section := range credentialsFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") {
			profiles = append(profiles, section.Name())
		}
	}

	configPath := ExpandHomeDirectory(*setCommand.Flags.ConfigFilePath)
	configFile, err := ini.Load(configPath)
	if err != nil {
		fmt.Printf("Fail to read AWS config file: %v", err)
		os.Exit(1)
	}

	for _, section := range configFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") &&
			section.HasKey("role_arn") &&
			section.HasKey("source_profile") {
			profiles = append(profiles, section.Name())
		}
	}

	return profiles
}

func (setCommand SetCommand) Handle(arguments []string) {
	command := setCommand.Command
	command.Parse(arguments)
	if command.Parsed() {
		fmt.Println(*setCommand.Flags.Pattern)
		profiles := getAWSProfiles(setCommand)
		joinedProfiles := strings.Join(profiles, "\n")

		fzfCommand := fmt.Sprintf("echo -e '%s' | fzf-tmux --height 30%% --reverse -1 -0 --header 'Select AWS profile' --query '%s'",
								joinedProfiles,
								*setCommand.Flags.Pattern)
		shellCommand := exec.Command("bash", "-c", fzfCommand)
		shellCommand.Stdin = os.Stdin
		shellCommand.Stderr = os.Stderr

		shellOutput, err := shellCommand.Output()
		if err != nil {
			// should only exit with code 0 when the error is caused by Ctrl+C
			// temporarily assume all the errors are caused by Ctrl+C for now
			os.Exit(0)
		}

		selectedProfile := string(shellOutput)
		fmt.Printf("%s", selectedProfile)
	}
}
