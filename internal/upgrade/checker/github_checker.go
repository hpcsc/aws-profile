package checker

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	latestReleaseUrl = "https://api.github.com/repos/hpcsc/aws-profile/releases/latest"
)

type githubLatestReleaseResponse struct {
	Assets []githubLatestReleaseResponseAsset `json:"assets"`
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

func (c githubChecker) LatestVersionUrl() (string, error) {
	bodyContent, err := c.getUrl(latestReleaseUrl)
	if err != nil {
		return "", err
	}

	unmarshalledResponse := githubLatestReleaseResponse{}
	err = json.Unmarshal(bodyContent, &unmarshalledResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %s", bodyContent)
	}

	for _, asset := range unmarshalledResponse.Assets {
		if strings.Contains(strings.ToLower(asset.Name), c.os) {
			return asset.Url, nil
		}
	}

	return "", fmt.Errorf("download url for os %s not found", c.os)
}
