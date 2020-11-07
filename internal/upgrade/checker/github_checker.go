package checker

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	latestReleaseUrl  = "https://api.github.com/repos/hpcsc/aws-profile/releases/latest"
	getCommitByTagUrl = "https://api.github.com/repos/hpcsc/aws-profile/commits/"
)

type githubGetCommitByTagResponse struct {
	Sha string `json:"sha"`
}

type githubLatestReleaseResponse struct {
	TagName string                             `json:"tag_name"`
	Assets  []githubLatestReleaseResponseAsset `json:"assets"`
}

type githubLatestReleaseResponseAsset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

type githubChecker struct {
	os     string
	getUrl func(string) ([]byte, error)
}

func newGithubChecker(os string, getUrl func(string) ([]byte, error)) checker {
	return githubChecker{os: os, getUrl: getUrl}
}

func (c githubChecker) LatestVersionUrl() (string, string, error) {
	bodyContent, err := c.getUrl(latestReleaseUrl)
	if err != nil {
		return "", "", err
	}

	unmarshalledResponse := githubLatestReleaseResponse{}
	err = json.Unmarshal(bodyContent, &unmarshalledResponse)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response body: %s", bodyContent)
	}

	for _, asset := range unmarshalledResponse.Assets {
		if strings.Contains(strings.ToLower(asset.Name), c.os) {
			return asset.Url, c.commitHashForTag(unmarshalledResponse.TagName), nil
		}
	}

	return "", "", fmt.Errorf("download url for os %s not found", c.os)
}

func (c githubChecker) commitHashForTag(tag string) string {
	// Purposely ignore error here.
	// Commit hash is only used to check whether users has latest version in their machines.
	// Any error here, just upgrade no matter what

	bodyContent, err := c.getUrl(getCommitByTagUrl + tag)
	if err != nil {
		return ""
	}

	unmarshalledResponse := githubGetCommitByTagResponse{}
	err = json.Unmarshal(bodyContent, &unmarshalledResponse)
	if err != nil {
		return ""
	}

	return unmarshalledResponse.Sha
}
