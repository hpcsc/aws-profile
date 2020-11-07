package checker

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var stubGithubGetUrl = func(url string) ([]byte, error) {
	response, err := ioutil.ReadFile("testdata/github-latest-release-response.json")
	if err != nil {
		fmt.Errorf("failed to read response file: %v", err)
		os.Exit(1)
	}

	return response, nil
}

func TestGithubChecker_LatestVersionUrl(t *testing.T) {
	var testCases = []struct {
		os          string
		expectedUrl string
	}{
		{"windows", "https://github.com/hpcsc/aws-profile/releases/download/v0.4.0/aws-profile-windows.exe"},
		{"linux", "https://github.com/hpcsc/aws-profile/releases/download/v0.4.0/aws-profile-linux"},
		{"macos", "https://github.com/hpcsc/aws-profile/releases/download/v0.4.0/aws-profile-macos"},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("return link to latest %s binary when os is %s", tt.os, tt.os), func(t *testing.T) {
			c := newGithubChecker(tt.os, stubGithubGetUrl)

			url, err := c.LatestVersionUrl()

			require.NoError(t, err)
			require.Equal(t, tt.expectedUrl, url)
		})
	}

	t.Run("return error when fail to get url", func(t *testing.T) {
		c := newGithubChecker("linux", func(url string) ([]byte, error) {
			return nil, errors.New("some error")
		})

		_, err := c.LatestVersionUrl()

		require.Error(t, err)
		require.Contains(t, err.Error(), "some error")
	})

	t.Run("return error when no asset for given os found", func(t *testing.T) {
		c := newGithubChecker("bsd", stubGithubGetUrl)

		_, err := c.LatestVersionUrl()

		require.Error(t, err)
		require.Contains(t, err.Error(), "download url for os bsd not found")
	})
}
