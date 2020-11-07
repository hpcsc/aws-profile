package upgrade

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/upgrade/checker"
	"github.com/hpcsc/aws-profile/internal/upgrade/httpclient"
	"os"
	"runtime"
)

func ToLatest(binary string, includePrerelease bool, currentVersion string) (string, error) {
	osName, err := getOSName(runtime.GOOS)
	if err != nil {
		return "", err
	}

	c := checker.NewChecker(osName, includePrerelease)
	url, version, err := c.LatestVersionUrl()
	if err != nil {
		return "", err
	}

	if currentVersion == version {
		return fmt.Sprintf("aws-profile is already at latest version (%s)", version), nil
	}

	newFileName := "./aws-profile.new"
	err = httpclient.DownloadFile(newFileName, url)
	if err != nil {
		return "", err
	}

	err = os.Chmod(newFileName, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to change downloaded file permission: %v", err)
	}

	old := binary + ".old"
	os.Remove(old)

	err = os.Rename(binary, old)
	if err != nil {
		return "", fmt.Errorf("failed to rename current executable: %v", err)
	}

	if err := os.Rename(newFileName, binary); err != nil {
		os.Rename(old, binary)
		return "", fmt.Errorf("failed to rename downloaded binary %s to %s: %v", newFileName, binary, err)
	}

	return fmt.Sprintf("aws-profile upgraded to latest version (%s)", version), nil
}

func getOSName(goos string) (string, error) {
	switch goos {
	case "windows":
		return "windows", nil
	case "linux":
		return "linux", nil
	case "darwin":
		return "macos", nil
	default:
		return "", fmt.Errorf("not supported os: %s", goos)
	}
}
