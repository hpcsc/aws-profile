package handlers

import (
	"github.com/hpcsc/aws-profile/internal/awsconfig"
)

type SelectProfileFn func(awsconfig.Profiles, string) ([]byte, error)
