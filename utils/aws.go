package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetAWSCredentials(profile *AWSProfile) (credentials.Value, error) {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	credentials := stscreds.NewCredentials(session, profile.RoleArn, func(p *stscreds.AssumeRoleProvider) {
		if profile.MFASerialNumber != "" {
			p.SerialNumber = aws.String(profile.MFASerialNumber)
			p.TokenProvider = stscreds.StdinTokenProvider
		}
		p.RoleSessionName = "aws-profile-session"
	})

	return credentials.Get()
}
