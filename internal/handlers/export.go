package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hpcsc/aws-profile/internal/aws"
	"github.com/hpcsc/aws-profile/internal/io"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

type GetAWSCredentialsFn func(*aws.AWSProfile, time.Duration) (credentials.Value, error)

type ExportHandler struct {
	SubCommand        *kingpin.CmdClause
	IsWindows         bool
	SelectProfile     SelectProfileFn
	GetAWSCredentials GetAWSCredentialsFn
	Arguments         ExportCommandArguments
}

type ExportCommandArguments struct {
	Pattern  *string
	Duration *string
}

func NewExportHandler(app *kingpin.Application, isWindows bool, selectProfileFn SelectProfileFn, getAWSCredentialsFn GetAWSCredentialsFn) ExportHandler {
	subCommand := app.Command("export", `print commands to set environment variables for assuming a AWS role

To execute the command without printing it to console:

- For Linux/MacOS, execute: "eval $(aws-profile export)"

- For Windows, execute: "Invoke-Expression (path\to\aws-profile.exe export)"`)

	pattern := subCommand.Arg("pattern", "Filter profiles by given pattern").String()
	duration := subCommand.Flag("duration", "AWS temporary session token duration. Example of valid duration: 5s, 15m").Short('d').Default("15m").String()

	return ExportHandler{
		SubCommand:        subCommand,
		IsWindows:         isWindows,
		SelectProfile:     selectProfileFn,
		GetAWSCredentials: getAWSCredentialsFn,
		Arguments: ExportCommandArguments{
			Pattern:  pattern,
			Duration: duration,
		},
	}
}

func (handler ExportHandler) Handle(globalArguments GlobalArguments) (bool, string) {
	configFile, readConfigErr := io.ReadFile(*globalArguments.ConfigFilePath)
	if readConfigErr != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", readConfigErr)
	}

	var duration, parseDurationErr = time.ParseDuration(*handler.Arguments.Duration)
	if parseDurationErr != nil {
		return false, parseDurationErr.Error()
	}

	if duration < time.Duration(15)*time.Minute {
		return false, "Minimum duration is 15 minutes"
	}

	processor := aws.AWSSharedCredentialsProcessor{
		CredentialsFile: ini.Empty(),
		ConfigFile:      configFile,
	}

	profiles := processor.GetProfilesFromCredentialsAndConfig()

	selectProfileResult, selectProfileErr := handler.SelectProfile(profiles, *handler.Arguments.Pattern)
	if selectProfileErr != nil {
		// cancel by user
		return true, ""
	}

	trimmedSelectedProfileResult := strings.TrimSuffix(string(selectProfileResult), "\n")
	profile := profiles.FindProfileInConfigFile(trimmedSelectedProfileResult)

	credentialsValue, getCredentialsErr := handler.GetAWSCredentials(profile, duration)
	if getCredentialsErr != nil {
		return false, getCredentialsErr.Error()
	}

	if handler.IsWindows {
		return true, formatOutputForWindows(credentialsValue, profile)
	} else {
		return true, formatOutputForLinuxAndMacOS(credentialsValue, profile)
	}
}

func formatOutputForWindows(credentialsValue credentials.Value, profile *aws.AWSProfile) string {
	output := fmt.Sprintf("$env:AWS_ACCESS_KEY_ID = '%s'; $env:AWS_SECRET_ACCESS_KEY = '%s'; $env:AWS_SESSION_TOKEN = '%s'",
		credentialsValue.AccessKeyID,
		credentialsValue.SecretAccessKey,
		credentialsValue.SessionToken,
	)

	if profile.Region == "" {
		return output
	}

	return fmt.Sprintf("%s; $env:AWS_REGION = '%s'; $env:AWS_DEFAULT_REGION = '%s'",
		output,
		profile.Region,
		profile.Region)
}

func formatOutputForLinuxAndMacOS(credentialsValue credentials.Value, profile *aws.AWSProfile) string {
	output := fmt.Sprintf("export AWS_ACCESS_KEY_ID='%s' AWS_SECRET_ACCESS_KEY='%s' AWS_SESSION_TOKEN='%s'",
		credentialsValue.AccessKeyID,
		credentialsValue.SecretAccessKey,
		credentialsValue.SessionToken,
	)

	if profile.Region == "" {
		return output
	}

	return fmt.Sprintf("%s AWS_REGION='%s' AWS_DEFAULT_REGION='%s'",
		output,
		profile.Region,
		profile.Region)
}
