package utils

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"gopkg.in/ini.v1"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func ReadFile(filePath string) (*ini.File, error) {
	path := ExpandHomeDirectory(filePath)
	return ini.Load(path)
}

func ExpandHomeDirectory(filePath string) string {
	usr, _ := user.Current()
	homeDirectory := usr.HomeDir
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(homeDirectory, filePath[2:])
	}

	return filePath
}

type GlobalArguments struct {
	CredentialsFilePath *string
	ConfigFilePath      *string
}

type Handler interface {
	Handle(globalArguments GlobalArguments) (bool, string)
}

type SelectProfileFn func(AWSProfiles, string) ([]byte, error)
type WriteToFileFn func(*ini.File, string)
type GetAWSCredentialsFn func(*AWSProfile, time.Duration) (credentials.Value, error)
type GetAWSCallerIdentityFn func() (string, error)
