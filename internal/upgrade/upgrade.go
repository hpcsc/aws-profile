package upgrade

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/upgrade/checker"
	"github.com/hpcsc/aws-profile/internal/upgrade/httpclient"
	"os"
	"path/filepath"
	"runtime"
)

func ToLatest(currentBinaryPath string, includePrerelease bool, currentVersion string) (string, error) {
	currentBinaryName := filepath.Base(currentBinaryPath)

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
		return fmt.Sprintf("%s is already at latest version (%s)", currentBinaryName, version), nil
	}

	newFileName := fmt.Sprintf("./%s.new", currentBinaryName)
	err = httpclient.DownloadFile(newFileName, url)
	if err != nil {
		return "", err
	}

	err = os.Chmod(newFileName, 0755) // #nosec
	if err != nil {
		return "", fmt.Errorf("failed to change downloaded file permission: %v", err)
	}

	old := currentBinaryPath + ".old"
	_ = os.Remove(old)

	err = os.Rename(currentBinaryPath, old)
	if err != nil {
		return "", fmt.Errorf("failed to rename current executable: %v", err)
	}

	if err := os.Rename(newFileName, currentBinaryPath); err != nil {
		renameErr := fmt.Errorf("failed to rename downloaded binary %s to %s: %v", newFileName, currentBinaryPath, err)

		if err := os.Rename(old, currentBinaryPath); err != nil {
			renameErr = fmt.Errorf("failed to recover original binary: %v, original error: %v", err, renameErr)
		}

		return "", renameErr
	}

	return fmt.Sprintf("%s upgraded to latest version (%s)", currentBinaryName, version), nil
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
