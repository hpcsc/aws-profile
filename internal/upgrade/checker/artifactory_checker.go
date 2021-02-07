package checker

import (
	"fmt"
	"strings"
)

const (
	artifactoryLatestVersionUrl = "https://hpcsc.jfrog.io/artifactory/aws-profile/latest-version"
)

type artifactoryChecker struct {
	os     string
	getUrl func(string) ([]byte, error)
}

func newArtifactoryChecker(os string, getUrl func(string) ([]byte, error)) checker {
	return artifactoryChecker{os: os, getUrl: getUrl}
}

func (c artifactoryChecker) LatestVersionUrl() (string, string, error) {
	bodyContent, err := c.getUrl(artifactoryLatestVersionUrl)
	if err != nil {
		return "", "", err
	}

	latestVersion := strings.TrimSpace(string(bodyContent))
	return fmt.Sprintf("https://hpcsc.jfrog.io/artifactory/aws-profile/%s/aws-profile-%s", latestVersion, strings.ToLower(c.os)), latestVersion, nil
}
