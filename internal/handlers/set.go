package handlers

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/config"
	"github.com/hpcsc/aws-profile/internal/io"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type SetHandler struct {
	SubCommand    *kingpin.CmdClause
	Arguments     SetCommandArguments
	SelectProfile SelectProfileFn
	WriteToFile   WriteToFileFn
}

type SetCommandArguments struct {
	Pattern *string
}

func NewSetHandler(app *kingpin.Application, selectProfileFn SelectProfileFn, writeToFileFn WriteToFileFn) SetHandler {
	subCommand := app.Command("set", "set default profile with credentials of selected profile")

	pattern := subCommand.Arg("pattern", "Filter profiles by given pattern").String()

	return SetHandler{
		SubCommand: subCommand,
		Arguments: SetCommandArguments{
			Pattern: pattern,
		},
		SelectProfile: selectProfileFn,
		WriteToFile:   writeToFileFn,
	}
}

func (handler SetHandler) Handle(globalArguments GlobalArguments) (bool, string) {
	credentialsFile, err := io.ReadFile(globalArguments.CredentialsFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS credentials file: %v", err)
	}

	configFile, err := io.ReadFile(globalArguments.ConfigFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", err)
	}

	profiles := config.LoadProfilesFromConfigAndCredentials(credentialsFile, configFile)

	selectProfileResult, err := handler.SelectProfile(profiles, *handler.Arguments.Pattern)
	if err != nil {
		// should only exit with code 0 when the error is caused by Ctrl+C
		// temporarily assume all the errors are caused by Ctrl+C for now
		return true, ""
	}

	trimmedSelectedProfileResult := strings.TrimSuffix(string(selectProfileResult), "\n")

	if profiles.FindProfileInCredentialsFile(trimmedSelectedProfileResult) != nil {
		config.SetSelectedProfileAsDefault(trimmedSelectedProfileResult, credentialsFile, configFile)

		handler.WriteToFile(credentialsFile, globalArguments.CredentialsFilePath)
		handler.WriteToFile(configFile, globalArguments.ConfigFilePath)

		return true, fmt.Sprintf("=== [%s] -> [default] (%s)", trimmedSelectedProfileResult, globalArguments.CredentialsFilePath)
	} else if assumedProfile := profiles.FindProfileInConfigFile(trimmedSelectedProfileResult); assumedProfile != nil {
		config.SetSelectedAssumedProfileAsDefault(assumedProfile.ProfileName, configFile)

		handler.WriteToFile(configFile, globalArguments.ConfigFilePath)

		return true, fmt.Sprintf("=== [%s] -> [default] (%s)", assumedProfile.ProfileName, globalArguments.ConfigFilePath)
	} else {
		return false, fmt.Sprintf("=== profile [%s] not found in either credentials or config file", trimmedSelectedProfileResult)
	}
}
