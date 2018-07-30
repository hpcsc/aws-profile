package handlers

import (
		"fmt"
	"os/exec"
	"os"
		"strings"
	"gopkg.in/ini.v1"
	"bytes"
	"io/ioutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

type SetHandler struct {
	SubCommand *kingpin.CmdClause
	Arguments  SetCommandArguments
}

type SetCommandArguments struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
	Pattern *string
}

func NewSetHandler(app *kingpin.Application) SetHandler {
	subCommand := app.Command("set", "set default profile with credentials of selected profile (this command assumes fzf is already setup)")

	credentialsFilePath := subCommand.Arg("credentials-path", "Path to AWS Credentials file").Default("~/.aws/credentials").String()
	configFilePath := subCommand.Arg("config-path", "Path to AWS Config file").Default("~/.aws/config").String()
	pattern := subCommand.Arg("pattern", "Start the fzf finder with the given query").String()

	return SetHandler{
		SubCommand: subCommand,
		Arguments:   SetCommandArguments{
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

func (handler SetHandler) Handle() {
	credentialsFile, err := ReadFile(*handler.Arguments.CredentialsFilePath)
	if err != nil {
		fmt.Printf("Fail to read AWS credentials file: %v", err)
		os.Exit(1)
	}

	configFile, err := ReadFile(*handler.Arguments.ConfigFilePath)
	if err != nil {
		fmt.Printf("Fail to read AWS config file: %v", err)
		os.Exit(1)
	}

	credentialsProfiles := getProfilesFromCredentialsFile(credentialsFile)
	configAssumedProfiles := getAssumedProfilesFromConfigFile(configFile)

	joinedProfiles := strings.Join(append(credentialsProfiles, configAssumedProfiles...), "\n")

	fzfCommand := fmt.Sprintf("echo -e '%s' | fzf-tmux --height 30%% --reverse -1 -0 --header 'Select AWS profile' --query '%s'",
							joinedProfiles,
							*handler.Arguments.Pattern)
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

		writeToFile(credentialsFile, *handler.Arguments.CredentialsFilePath)
		writeToFile(configFile, *handler.Arguments.ConfigFilePath)
	} else if containsProfile(configAssumedProfiles, selectedProfile) {
		fmt.Printf("=== assuming AWS profile [%s]", selectedProfile)
		selectedRoleArn := configFile.Section(selectedProfile).Key("role_arn").Value()
		selectedSourceProfile := configFile.Section(selectedProfile).Key("source_profile").Value()

		configFile.Section("default").Key("role_arn").SetValue(selectedRoleArn)
		configFile.Section("default").Key("source_profile").SetValue(selectedSourceProfile)

		writeToFile(configFile, *handler.Arguments.ConfigFilePath)
	} else {
		fmt.Printf("=== profile [%s] not found in either credentials or config file", selectedProfile)
	}
}
