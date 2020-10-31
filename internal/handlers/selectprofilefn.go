package handlers

import (
	"github.com/hpcsc/aws-profile/internal/config"
)

type SelectProfileFn func(config.Profiles, string) ([]byte, error)
