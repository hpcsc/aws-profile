package handlers

import (
	"flag"
	"fmt"
	"os/exec"
	"os"
		"strings"
	"gopkg.in/ini.v1"
	"bytes"
	"io/ioutil"
)

type SetHandler struct {
	FlagSet *flag.FlagSet
	Flags   SetCommandFlags
}

type SetCommandFlags struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
	Pattern *string
}

func NewSetHandler() SetHandler {
	flagSet := flag.NewFlagSet("set", flag.ExitOnError)

	credentialsFilePath := flagSet.String("credentials-path", "~/.aws/credentials", "Path to AWS Credentials file")
	configFilePath := flagSet.String("config-path", "~/.aws/config", "Path to AWS Config file")
	pattern := flagSet.String("pattern", "", "Start the fzf finder with the given query")

	return SetHandler{
		FlagSet: flagSet,
		Flags:   SetCommandFlags{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
			Pattern: pattern,
		},
	}
}

func getProfilesFromCredentialsFile(credentialsFile *ini.File) []string {
	var profiles []string

	for _, section := range credentialsFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") {
			profiles = append(profiles, section.Name())
		}
	}

	return profiles
}

func getAssumedProfilesFromConfigFile(configFile *ini.File) []string {
	var profiles []string

	for _, section := range configFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") &&
			section.HasKey("role_arn") &&
			section.HasKey("source_profile") {
			profiles = append(profiles, section.Name())
		}
	}

	return profiles
}

func containsProfile(profiles []string, selected string) bool {
	for _, profile := range profiles {
		if strings.EqualFold(profile, selected) {
			return true
		}
	}
	return false
}

func writeToFile(file *ini.File, unexpandedFilePath string) {
	var buffer bytes.Buffer
	_, err := file.WriteTo(&buffer)

	filePath := ExpandHomeDirectory(unexpandedFilePath)

	if err != nil {
		fmt.Printf("Fail to write to file %s: %v", filePath, err)
		os.Exit(1)
	}

	ioutil.WriteFile(filePath, buffer.Bytes(), 0600)
}

func (handler SetHandler) Handle(arguments []string) {
	flagSet := handler.FlagSet
	flagSet.Parse(arguments)
	if flagSet.Parsed() {
		credentialsFile, err := ReadFile(*handler.Flags.CredentialsFilePath)
		if err != nil {
			fmt.Printf("Fail to read AWS credentials file: %v", err)
			os.Exit(1)
		}

		configFile, err := ReadFile(*handler.Flags.ConfigFilePath)
		if err != nil {
			fmt.Printf("Fail to read AWS config file: %v", err)
			os.Exit(1)
		}

		credentialsProfiles := getProfilesFromCredentialsFile(credentialsFile)
		configAssumedProfiles := getAssumedProfilesFromConfigFile(configFile)

		joinedProfiles := strings.Join(append(credentialsProfiles, configAssumedProfiles...), "\n")

		fzfCommand := fmt.Sprintf("echo -e '%s' | fzf-tmux --height 30%% --reverse -1 -0 --header 'Select AWS profile' --query '%s'",
								joinedProfiles,
								*handler.Flags.Pattern)
		shellCommand := exec.Command("bash", "-c", fzfCommand)
		shellCommand.Stdin = os.Stdin
		shellCommand.Stderr = os.Stderr

		shellOutput, err := shellCommand.Output()
		if err != nil {
			// should only exit with code 0 when the error is caused by Ctrl+C
			// temporarily assume all the errors are caused by Ctrl+C for now
			os.Exit(0)
		}

		selectedProfile := strings.TrimSuffix(string(shellOutput), "\n")

		if containsProfile(credentialsProfiles, selectedProfile) {
			fmt.Printf("=== setting AWS profile [%s] as default profile", selectedProfile)
			selectedKeyId := credentialsFile.Section(selectedProfile).Key("aws_access_key_id").Value()
			selectedAccessKey := credentialsFile.Section(selectedProfile).Key("aws_secret_access_key").Value()

			credentialsFile.Section("default").Key("aws_access_key_id").SetValue(selectedKeyId)
			credentialsFile.Section("default").Key("aws_secret_access_key").SetValue(selectedAccessKey)
			configFile.Section("default").DeleteKey("role_arn")
			configFile.Section("default").DeleteKey("source_profile")

			writeToFile(credentialsFile, *handler.Flags.CredentialsFilePath)
			writeToFile(configFile, *handler.Flags.ConfigFilePath)
		} else if containsProfile(configAssumedProfiles, selectedProfile) {
			fmt.Printf("=== assuming AWS profile [%s]", selectedProfile)
			selectedRoleArn := configFile.Section(selectedProfile).Key("role_arn").Value()
			selectedSourceProfile := configFile.Section(selectedProfile).Key("source_profile").Value()

			configFile.Section("default").Key("role_arn").SetValue(selectedRoleArn)
			configFile.Section("default").Key("source_profile").SetValue(selectedSourceProfile)

			writeToFile(configFile, *handler.Flags.ConfigFilePath)
		} else {
			fmt.Printf("=== profile [%s] not found in either credentials or config file", selectedProfile)
		}
	}
}
