package handlers

import "github.com/hpcsc/aws-profile/internal/aws"

type SelectProfileFn func(aws.AWSProfiles, string) ([]byte, error)
