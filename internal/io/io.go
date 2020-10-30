package io

import (
	"bytes"
	"fmt"
	"github.com/hpcsc/aws-profile/internal/utils"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path"
)

func WriteToFile(file *ini.File, unexpandedFilePath string) {
	var buffer bytes.Buffer
	_, err := file.WriteTo(&buffer)

	filePath := utils.ExpandHomeDirectory(unexpandedFilePath)

	if err != nil {
		fmt.Printf("Fail to write to file %s: %v", filePath, err)
		os.Exit(1)
	}

	ioutil.WriteFile(filePath, buffer.Bytes(), 0600)
}

const awsProfileHome = "~/.aws-profile"
const cachedCallerIdentityFileName = "cached-caller-identity"

func createCachedCallerIdentityFileIfNotExists() (string, error) {
	expandedHomeDir := utils.ExpandHomeDirectory(awsProfileHome)

	if _, statError := os.Stat(expandedHomeDir); os.IsNotExist(statError) {
		makeDirError := os.Mkdir(expandedHomeDir, os.FileMode(0755))
		if makeDirError != nil {
			return "", makeDirError
		}
	}

	return path.Join(expandedHomeDir, cachedCallerIdentityFileName), nil
}

func ReadCachedCallerIdentity() (string, error) {
	cachedCallerIdentityFile, createError := createCachedCallerIdentityFileIfNotExists()
	if createError != nil {
		return "", createError
	}

	callerIdentity, readError := ioutil.ReadFile(cachedCallerIdentityFile)
	if readError != nil {
		return "", readError
	}

	return string(callerIdentity), nil
}

func WriteCachedCallerIdentity(callerIdentity string) error {
	cachedCallerIdentityFile, createError := createCachedCallerIdentityFileIfNotExists()
	if createError != nil {
		return createError
	}

	return ioutil.WriteFile(cachedCallerIdentityFile, []byte(callerIdentity), os.FileMode(0644))
}

func ReadFile(filePath string) (*ini.File, error) {
	path := utils.ExpandHomeDirectory(filePath)
	return ini.Load(path)
}
