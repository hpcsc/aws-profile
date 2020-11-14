package checker

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	bintrayMasterPackageUrl = "https://api.bintray.com/packages/hpcsc/aws-profile/master"
)

type bintrayMasterPackageResponse struct {
	LatestVersion string `json:"latest_version"`
}

type bintrayChecker struct {
	os     string
	getUrl func(string) ([]byte, error)
}

func newBintrayChecker(os string, getUrl func(string) ([]byte, error)) checker {
	return bintrayChecker{os: os, getUrl: getUrl}
}

func (c bintrayChecker) LatestVersionUrl() (string, string, error) {
	bodyContent, err := c.getUrl(bintrayMasterPackageUrl)
	if err != nil {
		return "", "", err
	}

	var unmarshalledResponse bintrayMasterPackageResponse
	err = json.Unmarshal(bodyContent, &unmarshalledResponse)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response body: %s", bodyContent)
	}

	latestVersion := unmarshalledResponse.LatestVersion
	return fmt.Sprintf("https://dl.bintray.com/hpcsc/aws-profile/aws-profile-%s-%s", strings.ToLower(c.os), latestVersion), latestVersion, nil
}
