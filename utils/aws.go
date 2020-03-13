package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"time"
)

func GetAWSCredentials(profile *AWSProfile, duration time.Duration) (credentials.Value, error) {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	credentials := stscreds.NewCredentials(session, profile.RoleArn, func(p *stscreds.AssumeRoleProvider) {
		if profile.MFASerialNumber != "" {
			p.SerialNumber = aws.String(profile.MFASerialNumber)
			p.TokenProvider = stscreds.StdinTokenProvider
		}
		p.RoleSessionName = fmt.Sprintf("aws-profile-%d", time.Now().UnixNano())
		p.Duration = duration
	})

	return credentials.Get()
}
