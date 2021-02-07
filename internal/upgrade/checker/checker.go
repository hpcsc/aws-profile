package checker

import (
	"github.com/hpcsc/aws-profile/internal/upgrade/httpclient"
)

type checker interface {
	LatestVersionUrl() (string, string, error)
}

func NewChecker(os string, includePrerelease bool) checker {
	if includePrerelease {
		return newArtifactoryChecker(os, httpclient.GetUrl)
	}

	return newGithubChecker(os, httpclient.GetUrlWithAuthorization)
}
