package handlers

import (
	"errors"
	"fmt"
	"github.com/hpcsc/aws-profile/internal/awsconfig"
	"github.com/hpcsc/aws-profile/internal/config"
	"github.com/hpcsc/aws-profile/internal/io"
	"github.com/hpcsc/aws-profile/internal/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type SetHandler struct {
	SubCommand    *kingpin.CmdClause
	Arguments     SetCommandArguments
	SelectProfile SelectProfileFn
	WriteToFile   WriteToFileFn
	Config        *config.Config
}

type SetCommandArguments struct {
	Pattern *string
}

func NewSetHandler(app *kingpin.Application, config *config.Config, selectProfileFn SelectProfileFn, writeToFileFn WriteToFileFn) SetHandler {
	subCommand := app.Command("set", "set default profile with credentials of selected profile")

	pattern := subCommand.Arg("pattern", "Filter profiles by given pattern").String()

	return SetHandler{
		SubCommand: subCommand,
		Arguments: SetCommandArguments{
			Pattern: pattern,
		},
		SelectProfile: selectProfileFn,
		WriteToFile:   writeToFileFn,
		Config:        config,
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

	profiles := awsconfig.LoadProfilesFromConfigAndCredentials(credentialsFile, configFile)

	selectProfileResult, err := handler.SelectProfile(profiles, *handler.Arguments.Pattern, handler.Config)
	var cancelled *utils.CancelledError
	if errors.As(err, &cancelled) {
		return true, ""
	}

	if err != nil {
		return false, fmt.Sprintf("Failed to select profile: %v", err)
	}

	trimmedSelectedProfileResult := strings.TrimSuffix(string(selectProfileResult), "\n")

	if profiles.FindProfileInCredentialsFile(trimmedSelectedProfileResult) != nil {
		awsconfig.SetSelectedProfileAsDefault(trimmedSelectedProfileResult, credentialsFile, configFile)

		if err := handler.WriteToFile(credentialsFile, globalArguments.CredentialsFilePath); err != nil {
			return false, err.Error()
		}

		if err := handler.WriteToFile(configFile, globalArguments.ConfigFilePath); err != nil {
			return false, err.Error()
		}

		return true, fmt.Sprintf("=== [%s] -> [default] (%s)", trimmedSelectedProfileResult, globalArguments.CredentialsFilePath)
	} else if assumedProfile := profiles.FindProfileInConfigFile(trimmedSelectedProfileResult); assumedProfile != nil {
		awsconfig.SetSelectedAssumedProfileAsDefault(assumedProfile.ProfileName, configFile)

		if err := handler.WriteToFile(configFile, globalArguments.ConfigFilePath); err != nil {
			return false, err.Error()
		}

		return true, fmt.Sprintf("=== [%s] -> [default] (%s)", assumedProfile.ProfileName, globalArguments.ConfigFilePath)
	} else {
		return false, fmt.Sprintf("=== profile [%s] not found in either credentials or config file", trimmedSelectedProfileResult)
	}
}
