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
	response, err := ioutil.ReadFile("testdata/bintray-get-package-response.json")
	if err != nil {
		fmt.Printf("failed to read response file: %v", err)
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
		{"windows", "https://dl.bintray.com/hpcsc/aws-profile/aws-profile-windows-9c1ab3ae40ff6697a567c272a47337b1044506e3", "9c1ab3ae40ff6697a567c272a47337b1044506e3"},
		{"linux", "https://dl.bintray.com/hpcsc/aws-profile/aws-profile-linux-9c1ab3ae40ff6697a567c272a47337b1044506e3", "9c1ab3ae40ff6697a567c272a47337b1044506e3"},
		{"macos", "https://dl.bintray.com/hpcsc/aws-profile/aws-profile-macos-9c1ab3ae40ff6697a567c272a47337b1044506e3", "9c1ab3ae40ff6697a567c272a47337b1044506e3"},
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
}
