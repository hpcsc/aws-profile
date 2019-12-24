package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ExportHandler struct {
	SubCommand *kingpin.CmdClause
	Arguments  ExportCommandArguments
}

type ExportCommandArguments struct {
	roleToAssume *string
}

func NewExportHandler(app *kingpin.Application) ExportHandler {
	subCommand := app.Command("export", "print commands to set environment variables for assuming a AWS role")

	roleToAssume := subCommand.Arg("role-to-assume", "AWS role to assume").Required().String()

	return ExportHandler {
		SubCommand: subCommand,
		Arguments:   ExportCommandArguments{
			roleToAssume: roleToAssume,
		},
	}
}

func (handler ExportHandler) Handle() (bool, string) {
	mfaSerialNumber := "arn:aws:iam::697469898979:mfa/desktop-cli"

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	credentials := stscreds.NewCredentials(session, *handler.Arguments.roleToAssume, func(p *stscreds.AssumeRoleProvider) {
		p.SerialNumber = aws.String(mfaSerialNumber)
		p.TokenProvider = stscreds.StdinTokenProvider
		p.RoleSessionName = "aws-profile-utils-session"
	})

	credentialsValue, err := credentials.Get()
	if err != nil {
		return false, err.Error()
	}

	linuxExport := fmt.Sprintf("export AWS_ACCESS_KEY_ID='%s' AWS_SECRET_ACCESS_KEY='%s' AWS_SESSION_TOKEN='%s'\n",
		credentialsValue.AccessKeyID,
		credentialsValue.SecretAccessKey,
		credentialsValue.SessionToken)

	windowsExport := fmt.Sprintf("$env:AWS_ACCESS_KEY_ID = '%s'; $env:AWS_SECRET_ACCESS_KEY = '%s'; $env:AWS_SESSION_TOKEN = '%s'\n",
		credentialsValue.AccessKeyID,
		credentialsValue.SecretAccessKey,
		credentialsValue.SessionToken)

	output := fmt.Sprintf(`LINUX or MACOS
================================
Execute the following command in your shell:

%s

To unset those environment variables: 

unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN\n
								
WINDOWS
================================
Execute the following command in Powershell: 

%s

To unset those environment variables: 

Remove-Item Env:\AWS_ACCESS_KEY_ID; Remove-Item Env:\AWS_SECRET_ACCESS_KEY; Remove-Item Env:\AWS_SESSION_TOKEN`,
		linuxExport,
		windowsExport)
	return true, output
}
