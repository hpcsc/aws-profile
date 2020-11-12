package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hpcsc/aws-profile/internal/awsconfig"
	"strings"
	"time"
)

func GetAWSCredentials(profile *awsconfig.Profile, duration time.Duration) (credentials.Value, error) {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile.SourceProfile,
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

func GetAWSCallerIdentity() (string, error) {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	stsClient := sts.New(session)

	output, error := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if error != nil {
		return "", error
	}

	splittedArn := strings.Split(*output.Arn, "/")
	if len(splittedArn) < 2 {
		return *output.Arn, nil
	}

	roleName := splittedArn[1]
	return fmt.Sprintf("role %s@%s (env)", roleName, *output.Account), nil
}
