package utils

import (
	"bytes"
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
)

func WriteToFile(file *ini.File, unexpandedFilePath string) {
	var buffer bytes.Buffer
	_, err := file.WriteTo(&buffer)

	filePath := ExpandHomeDirectory(unexpandedFilePath)

	if err != nil {
		fmt.Printf("Fail to write to file %s: %v", filePath, err)
		os.Exit(1)
	}

	ioutil.WriteFile(filePath, buffer.Bytes(), 0600)
}

const cachedCallerIdentityFile = "~/.aws-profile-cached-caller-identity"

func ReadCachedCallerIdentity() (string, error) {
	callerIdentity, error := ioutil.ReadFile(ExpandHomeDirectory(cachedCallerIdentityFile))
	if error != nil {
		return "", error
	}

	return string(callerIdentity), nil
}

func WriteCachedCallerIdentity(callerIdentity string) error {
	return ioutil.WriteFile(ExpandHomeDirectory(cachedCallerIdentityFile), []byte(callerIdentity), 0600)
}
