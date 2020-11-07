package checker

import (
	"github.com/hpcsc/aws-profile/internal/upgrade/httpclient"
)

type checker interface {
	LatestVersionUrl() (string, string, error)
}

func NewChecker(os string, includePrerelease bool) checker {
	if includePrerelease {
		return newBintrayChecker(os, httpclient.GetUrl)
	}

	return newGithubChecker(os, httpclient.GetUrl)
}
