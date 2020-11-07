package checker

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var stubBintrayGetUrl = func(url string) ([]byte, error) {
	response, err := ioutil.ReadFile("testdata/bintray-files-response.json")
	if err != nil {
		fmt.Errorf("failed to read response file: %v", err)
		os.Exit(1)
	}

	return response, nil
}

func TestBintrayChecker_LatestVersionUrl(t *testing.T) {
	var testCases = []struct {
		os              string
		expectedUrl     string
		expectedVersion string
	}{
		{"windows", "https://dl.bintray.com/hpcsc/aws-profile/aws-profile-Windows-ebcda1baa76b902ecc035b5e5a232a488aa66cb0", "ebcda1baa76b902ecc035b5e5a232a488aa66cb0"},
		{"linux", "https://dl.bintray.com/hpcsc/aws-profile/aws-profile-Linux-ebcda1baa76b902ecc035b5e5a232a488aa66cb0", "ebcda1baa76b902ecc035b5e5a232a488aa66cb0"},
		{"macos", "https://dl.bintray.com/hpcsc/aws-profile/aws-profile-macOS-ebcda1baa76b902ecc035b5e5a232a488aa66cb0", "ebcda1baa76b902ecc035b5e5a232a488aa66cb0"},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("return version and link to latest %s binary when os is %s", tt.os, tt.os), func(t *testing.T) {
			c := newBintrayChecker(tt.os, stubBintrayGetUrl)

			url, version, err := c.LatestVersionUrl()

			require.NoError(t, err)
			require.Equal(t, tt.expectedUrl, url)
			require.Equal(t, tt.expectedVersion, version)
		})
	}

	t.Run("return error when fail to get url", func(t *testing.T) {
		c := newBintrayChecker("linux", func(url string) ([]byte, error) {
			return nil, errors.New("some error")
		})

		_, _, err := c.LatestVersionUrl()

		require.Error(t, err)
		require.Contains(t, err.Error(), "some error")
	})

	t.Run("return error when no asset for given os found", func(t *testing.T) {
		c := newBintrayChecker("bsd", stubBintrayGetUrl)

		_, _, err := c.LatestVersionUrl()

		require.Error(t, err)
		require.Contains(t, err.Error(), "download url for os bsd not found")
	})
}
