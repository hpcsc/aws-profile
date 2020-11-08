package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ExpandHomeDirectory(filePath string) string {
	if strings.HasPrefix(filePath, "~/") {
		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			// use log.Fatalf() instead of logger here because this function can be used by logger during init
			log.Fatalf("failed to get user home directory: %v", err)
		}

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
