package handlers

import (
	"fmt"
	"github.com/hpcsc/aws-profile/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type SetHandler struct {
	SubCommand *kingpin.CmdClause
	Arguments  SetCommandArguments
	SelectProfile utils.SelectProfileFn
	WriteToFile utils.WriteToFileFn
}

type SetCommandArguments struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
	Pattern *string
}

func NewSetHandler(app *kingpin.Application, selectProfileFn utils.SelectProfileFn, writeToFileFn utils.WriteToFileFn) SetHandler {
	subCommand := app.Command("set", "set default profile with credentials of selected profile")

	credentialsFilePath := subCommand.Flag("credentials-path", "Path to AWS Credentials file").Default("~/.aws/credentials").String()
	configFilePath := subCommand.Flag("config-path", "Path to AWS Config file").Default("~/.aws/config").String()
	pattern := subCommand.Arg("pattern", "Filter profiles by given pattern").String()

	return SetHandler{
		SubCommand: subCommand,
		Arguments:   SetCommandArguments{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
			Pattern: pattern,
		},
		SelectProfile: selectProfileFn,
		WriteToFile: writeToFileFn,
	}
}


func (handler SetHandler) Handle() (bool, string) {
	credentialsFile, err := utils.ReadFile(*handler.Arguments.CredentialsFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS credentials file: %v", err)
	}

	configFile, err := utils.ReadFile(*handler.Arguments.ConfigFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", err)
	}

	processor := utils.AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile: configFile,
	}

	profiles := processor.GetProfilesFromCredentialsAndConfig()

	selectProfileResult, err := handler.SelectProfile(profiles, *handler.Arguments.Pattern)
	if err != nil {
		// should only exit with code 0 when the error is caused by Ctrl+C
		// temporarily assume all the errors are caused by Ctrl+C for now
		return true, ""
	}

	trimmedSelectedProfileResult := strings.TrimSuffix(string(selectProfileResult), "\n")

	if profiles.FindProfileInCredentialsFile(trimmedSelectedProfileResult) != nil {
		processor.SetSelectedProfileAsDefault(trimmedSelectedProfileResult)

		handler.WriteToFile(processor.CredentialsFile, *handler.Arguments.CredentialsFilePath)
		handler.WriteToFile(processor.ConfigFile, *handler.Arguments.ConfigFilePath)

		return true, fmt.Sprintf("=== profile [default] in [%s] is set with credentials from profile [%s]", *handler.Arguments.CredentialsFilePath, trimmedSelectedProfileResult)
	} else if assumedProfile := profiles.FindProfileInConfigFile(trimmedSelectedProfileResult); assumedProfile != nil {
		processor.SetSelectedAssumedProfileAsDefault(assumedProfile.ProfileName)

		handler.WriteToFile(processor.ConfigFile, *handler.Arguments.ConfigFilePath)

		return true, fmt.Sprintf("=== profile [default] config in [%s] is set with configs from assumed [%s]", *handler.Arguments.ConfigFilePath, assumedProfile.ProfileName)
	} else {
		return false, fmt.Sprintf("=== profile [%s] not found in either credentials or config file", trimmedSelectedProfileResult)
	}
}
