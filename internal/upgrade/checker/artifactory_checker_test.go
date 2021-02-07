package checker

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var stubArtifactoryGetUrl = func(url string) ([]byte, error) {
	response, err := ioutil.ReadFile("testdata/artifactory-latest-version-response")
	if err != nil {
		fmt.Printf("failed to read response file: %v", err)
		os.Exit(1)
	}

	return response, nil
}

func TestArtifactoryChecker_LatestVersionUrl(t *testing.T) {
	var testCases = []struct {
		os              string
		expectedUrl     string
		expectedVersion string
	}{
		{"windows", "https://hpcsc.jfrog.io/artifactory/aws-profile/123/aws-profile-windows", "123"},
		{"linux", "https://hpcsc.jfrog.io/artifactory/aws-profile/123/aws-profile-linux", "123"},
		{"macos", "https://hpcsc.jfrog.io/artifactory/aws-profile/123/aws-profile-macos", "123"},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("return version and link to latest %s binary when os is %s", tt.os, tt.os), func(t *testing.T) {
			c := newArtifactoryChecker(tt.os, stubArtifactoryGetUrl)

			url, version, err := c.LatestVersionUrl()

			require.NoError(t, err)
			require.Equal(t, tt.expectedUrl, url)
			require.Equal(t, tt.expectedVersion, version)
		})
	}

	t.Run("return error when fail to get url", func(t *testing.T) {
		c := newArtifactoryChecker("linux", func(url string) ([]byte, error) {
			return nil, errors.New("some error")
		})

		_, _, err := c.LatestVersionUrl()

		require.Error(t, err)
		require.Contains(t, err.Error(), "some error")
	})
}
