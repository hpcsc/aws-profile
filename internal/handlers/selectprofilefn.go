package handlers

import (
	"github.com/hpcsc/aws-profile/internal/awsconfig"
	"github.com/hpcsc/aws-profile/internal/config"
)

type SelectProfileFn func(awsconfig.Profiles, string, *config.Config) ([]byte, error)
