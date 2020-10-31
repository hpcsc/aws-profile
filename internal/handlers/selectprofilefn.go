package handlers

import (
	"github.com/hpcsc/aws-profile/internal/config"
)

type SelectProfileFn func(config.AWSProfiles, string) ([]byte, error)
