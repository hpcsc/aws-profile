package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func ExpandHomeDirectory(filePath string) string {
	usr, _ := user.Current()
	homeDirectory := usr.HomeDir
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(homeDirectory, filePath[2:])
	}

	return filePath
}

func GetEnvVariableOrDefault(variableName string, defaultValue string) string {
	if value, exists := os.LookupEnv(variableName); exists {
		return value
	}

	return defaultValue
}
