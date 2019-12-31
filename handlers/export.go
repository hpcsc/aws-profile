package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hpcsc/aws-profile/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
	"strings"
)

type ExportHandler struct {
	SubCommand *kingpin.CmdClause
	IsWindows bool
	SelectProfile utils.SelectProfileFn
	GetAWSCredentials utils.GetAWSCredentialsFn
	Arguments  ExportCommandArguments
}

type ExportCommandArguments struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
	Pattern *string
}

func NewExportHandler(app *kingpin.Application, isWindows bool, selectProfileFn utils.SelectProfileFn, getAWSCredentialsFn utils.GetAWSCredentialsFn) ExportHandler {
	subCommand := app.Command("export", "print commands to set environment variables for assuming a AWS role")

	credentialsFilePath := subCommand.Flag("credentials-path", "Path to AWS Credentials file").Default("~/.aws/credentials").String()
	configFilePath := subCommand.Flag("config-path", "Path to AWS Config file").Default("~/.aws/config").String()
	pattern := subCommand.Arg("pattern", "Filter profiles by given pattern").String()

	return ExportHandler {
		SubCommand: subCommand,
		IsWindows: isWindows,
		SelectProfile: selectProfileFn,
		GetAWSCredentials: getAWSCredentialsFn,
		Arguments:   ExportCommandArguments{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
			Pattern: pattern,
		},
	}
}

func (handler ExportHandler) Handle() (bool, string) {
	configFile, readConfigErr := utils.ReadFile(*handler.Arguments.ConfigFilePath)
	if readConfigErr != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", readConfigErr)
	}

	processor := utils.AWSSharedCredentialsProcessor{
		CredentialsFile: ini.Empty(),
		ConfigFile: configFile,
	}

	profiles := processor.GetProfilesFromCredentialsAndConfig()

	selectProfileResult, selectProfileErr := handler.SelectProfile(profiles, *handler.Arguments.Pattern)
	if selectProfileErr != nil {
		// cancel by user
		return true, ""
	}

	trimmedSelectedProfileResult := strings.TrimSuffix(string(selectProfileResult), "\n")
	profile := profiles.FindProfileInConfigFile(trimmedSelectedProfileResult)

	credentialsValue, getCredentialsErr := handler.GetAWSCredentials(profile)
	if getCredentialsErr != nil {
		return false, getCredentialsErr.Error()
	}

	output := formatOutputByPlatform(handler.IsWindows, credentialsValue)
	return true, output
}

func formatOutputByPlatform(isWindows bool, credentialsValue credentials.Value) string {
	if isWindows {
		return fmt.Sprintf("$env:AWS_ACCESS_KEY_ID = '%s'; $env:AWS_SECRET_ACCESS_KEY = '%s'; $env:AWS_SESSION_TOKEN = '%s'",
			credentialsValue.AccessKeyID,
			credentialsValue.SecretAccessKey,
			credentialsValue.SessionToken)
	}

	return fmt.Sprintf("export AWS_ACCESS_KEY_ID='%s' AWS_SECRET_ACCESS_KEY='%s' AWS_SESSION_TOKEN='%s'",
		credentialsValue.AccessKeyID,
		credentialsValue.SecretAccessKey,
		credentialsValue.SessionToken)
}
