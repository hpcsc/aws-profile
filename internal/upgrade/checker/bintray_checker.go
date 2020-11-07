package checker

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	bintrayFilesUrl = "https://api.bintray.com/packages/hpcsc/aws-profile/master/files"
)

type bintrayFilesResponsestruct struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Created string `json:"created"`
}

type bintrayChecker struct {
	os     string
	getUrl func(string) ([]byte, error)
}

func newBintrayChecker(os string, getUrl func(string) ([]byte, error)) checker {
	return bintrayChecker{os: os, getUrl: getUrl}
}

func (c bintrayChecker) LatestVersionUrl() (string, string, error) {
	bodyContent, err := c.getUrl(bintrayFilesUrl)
	if err != nil {
		return "", "", err
	}

	var unmarshalledResponse []bintrayFilesResponsestruct
	err = json.Unmarshal(bodyContent, &unmarshalledResponse)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response body: %s", bodyContent)
	}

	latestPath := ""
	latestVersion := ""
	latestCreated := time.Time{}
	for _, file := range unmarshalledResponse {
		created, _ := time.Parse(time.RFC3339, file.Created)
		if strings.Contains(strings.ToLower(file.Name), c.os) &&
			created.After(latestCreated) {
			latestPath = file.Path
			latestVersion = file.Version
			latestCreated = created
		}
	}

	if latestPath != "" {
		return fmt.Sprintf("https://dl.bintray.com/hpcsc/aws-profile/%s", latestPath), latestVersion, nil
	}

	return "", "", fmt.Errorf("download url for os %s not found", c.os)
}
