package handlers

import (
	"bytes"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type SetHandler struct {
	SubCommand *kingpin.CmdClause
	Arguments  SetCommandArguments
	SelectProfile SelectProfileFn
	WriteToFile WriteToFileFn
}

type SetCommandArguments struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
	Pattern *string
}

type SelectProfileFn func([]string, string) ([]byte, error)
type WriteToFileFn func(*ini.File, string)

func NewSetHandler(app *kingpin.Application, selectProfileFn SelectProfileFn, writeToFileFn WriteToFileFn) SetHandler {
	subCommand := app.Command("set", "set default profile with credentials of selected profile (this command assumes fzf is already setup)")

	credentialsFilePath := subCommand.Flag("credentials-path", "Path to AWS Credentials file").Default("~/.aws/credentials").String()
	configFilePath := subCommand.Flag("config-path", "Path to AWS Config file").Default("~/.aws/config").String()
	pattern := subCommand.Arg("pattern", "Start the fzf finder with the given query").String()

	finalSelectProfileFn := selectProfileByFzf
	if selectProfileFn != nil {
		finalSelectProfileFn = selectProfileFn
	}

	finalWriteToFileFn := writeToFile
	if writeToFileFn != nil {
		finalWriteToFileFn = writeToFileFn
	}

	return SetHandler{
		SubCommand: subCommand,
		Arguments:   SetCommandArguments{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
			Pattern: pattern,
		},
		SelectProfile: finalSelectProfileFn,
		WriteToFile: finalWriteToFileFn,
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
			profiles = append(profiles, fmt.Sprintf("assume %s:%s", section.Name(), section.Name()))
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

func getProfileSectionNameFromColonDelimited(selected string) string {
	if ! strings.Contains(selected, ":") {
		return selected
	}

	elements := strings.Split(selected, ":")
	return elements[1]
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

func selectProfileByFzf(combinedProfiles []string, pattern string) ([]byte, error) {
	joinedProfiles := strings.Join(combinedProfiles, "\n")

	fzfCommand := fmt.Sprintf("echo -e '%s' | fzf-tmux --height 30%% --reverse -1 -0 --with-nth=1 --delimiter=: --header 'Select AWS profile' --query '%s'",
		joinedProfiles,
		pattern)
	shellCommand := exec.Command("bash", "-c", fzfCommand)
	shellCommand.Stdin = os.Stdin
	shellCommand.Stderr = os.Stderr
	return shellCommand.Output()
}

func (handler SetHandler) Handle() (bool, string) {
	credentialsFile, err := ReadFile(*handler.Arguments.CredentialsFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS credentials file: %v", err)
	}

	configFile, err := ReadFile(*handler.Arguments.ConfigFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", err)
	}

	credentialsProfiles := getProfilesFromCredentialsFile(credentialsFile)
	configAssumedProfiles := getAssumedProfilesFromConfigFile(configFile)

	combinedProfiles := append(credentialsProfiles, configAssumedProfiles...)
	shellOutput, err := handler.SelectProfile(combinedProfiles, *handler.Arguments.Pattern)
	if err != nil {
		// should only exit with code 0 when the error is caused by Ctrl+C
		// temporarily assume all the errors are caused by Ctrl+C for now
		return true, ""
	}

	selectedValue := strings.TrimSuffix(string(shellOutput), "\n")

	if containsProfile(credentialsProfiles, selectedValue) {
		selectedKeyId := credentialsFile.Section(selectedValue).Key("aws_access_key_id").Value()
		selectedAccessKey := credentialsFile.Section(selectedValue).Key("aws_secret_access_key").Value()

		credentialsFile.Section("default").Key("aws_access_key_id").SetValue(selectedKeyId)
		credentialsFile.Section("default").Key("aws_secret_access_key").SetValue(selectedAccessKey)
		configFile.Section("default").DeleteKey("role_arn")
		configFile.Section("default").DeleteKey("source_profile")

		handler.WriteToFile(credentialsFile, *handler.Arguments.CredentialsFilePath)
		handler.WriteToFile(configFile, *handler.Arguments.ConfigFilePath)

		return true, fmt.Sprintf("=== profile [default] in [%s] is set with credentials from profile [%s]", *handler.Arguments.CredentialsFilePath, selectedValue)
	} else if containsProfile(configAssumedProfiles, selectedValue) {
		selectedProfile := getProfileSectionNameFromColonDelimited(selectedValue)

		selectedRoleArn := configFile.Section(selectedProfile).Key("role_arn").Value()
		selectedSourceProfile := configFile.Section(selectedProfile).Key("source_profile").Value()

		configFile.Section("default").Key("role_arn").SetValue(selectedRoleArn)
		configFile.Section("default").Key("source_profile").SetValue(selectedSourceProfile)

		handler.WriteToFile(configFile, *handler.Arguments.ConfigFilePath)

		return true, fmt.Sprintf("=== profile [default] config in [%s] is set with configs from assumed [%s]", *handler.Arguments.ConfigFilePath, selectedProfile)
	} else {
		return false, fmt.Sprintf("=== profile [%s] not found in either credentials or config file", selectedValue)
	}
}
