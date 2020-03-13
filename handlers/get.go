package handlers

import (
	"fmt"
	"github.com/hpcsc/aws-profile/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

type GetHandler struct {
	SubCommand *kingpin.CmdClause
	GetAWSCallerIdentityFn utils.GetAWSCallerIdentityFn
}

func NewGetHandler(app *kingpin.Application, getAWSCallerIdentityFn utils.GetAWSCallerIdentityFn) GetHandler {
	subCommand := app.Command("get", "get current AWS profile (that is set to default profile)")

	return GetHandler{
		SubCommand: subCommand,
		GetAWSCallerIdentityFn: getAWSCallerIdentityFn,
	}
}

func awsCredentialsEnvironmentVariablesSet() bool {
	var _, accessKeyIdExists = os.LookupEnv("AWS_ACCESS_KEY_ID")
	var _, secretAccessKeyExists = os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	var _, sessionTokenExists = os.LookupEnv("AWS_SESSION_TOKEN")

	return accessKeyIdExists && secretAccessKeyExists && sessionTokenExists
}

func (handler GetHandler) Handle(globalArguments utils.GlobalArguments) (bool, string) {
	if awsCredentialsEnvironmentVariablesSet() {
		var callerIdentityProfile, getCallerIdentityErr = handler.GetAWSCallerIdentityFn()
		if getCallerIdentityErr != nil {
			// purposely ignore error from AWS STS
			return true, "unknown"
		}

		return true, callerIdentityProfile
	}

	configFile, err := utils.ReadFile(*globalArguments.ConfigFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", err)
	}

	configDefaultSection, err := configFile.GetSection("default")
	if err == nil &&
		configDefaultSection.HasKey("role_arn") &&
		configDefaultSection.HasKey("source_profile") {

		defaultRoleArn := configDefaultSection.Key("role_arn").Value()
		defaultSourceProfile := configDefaultSection.Key("source_profile").Value()

		for _, section := range configFile.Sections() {
			if strings.Compare(section.Name(), "default") != 0 &&
				section.HasKey("role_arn") &&
				section.HasKey("source_profile") &&
				strings.Compare(section.Key("role_arn").Value(), defaultRoleArn) == 0 &&
				strings.Compare(section.Key("source_profile").Value(), defaultSourceProfile) == 0 {
				return true, section.Name()
			}
		}
	}

	credentialsFile, err := utils.ReadFile(*globalArguments.CredentialsFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS credentials file: %v", err)
	}

	credentialsDefaultSection, err := credentialsFile.GetSection("default")
	if err == nil &&
		credentialsDefaultSection.HasKey("aws_access_key_id") {

		defaultAWSAccessKeyId := credentialsDefaultSection.Key("aws_access_key_id").Value()

		for _, section := range credentialsFile.Sections() {
			if strings.Compare(section.Name(), "default") != 0 &&
				section.HasKey("aws_access_key_id") &&
				strings.Compare(section.Key("aws_access_key_id").Value(), defaultAWSAccessKeyId) == 0 {
				return true, fmt.Sprintf("%s\n", section.Name())
			}
		}
	}

	return true, ""
}
