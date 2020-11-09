package handlers

import (
	"fmt"
	"strings"

	"github.com/hpcsc/aws-profile/internal/config"
	"github.com/hpcsc/aws-profile/internal/io"
	"gopkg.in/alecthomas/kingpin.v2"
)

var Regions = []string{"ap-southeast-2", "us-west-1"}

type SelectRegionFn func([]string, string) ([]byte, error)

type SetRegionHandler struct {
	SubCommand   *kingpin.CmdClause
	SelectRegion SelectRegionFn
	WriteToFile  WriteToFileFn
}

type SetRegionCommandArguments struct {
	Pattern *string
}

func NewSetRegionHandler(app *kingpin.Application, selectRegionFn SelectRegionFn, writeToFileFn WriteToFileFn) SetRegionHandler {
	subCommand := app.Command("set-region", "set the region of the default profile")

	return SetRegionHandler{
		SubCommand:   subCommand,
		SelectRegion: selectRegionFn,
		WriteToFile:  writeToFileFn,
	}
}

func (handler SetRegionHandler) Handle(globalArguments GlobalArguments) (bool, string) {
	configFile, err := io.ReadFile(globalArguments.ConfigFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", err)
	}

	selectRegionResult, err := handler.SelectRegion(Regions, "Select an AWS region")
	if err != nil {
		// should only exit with code 0 when the error is caused by Ctrl+C
		// temporarily assume all the errors are caused by Ctrl+C for now
		return true, ""
	}
	trimmedSelectedRegionResult := strings.TrimSuffix(string(selectRegionResult), "\n")

	config.SetSelectedRegionAsDefault(trimmedSelectedRegionResult, configFile)
	handler.WriteToFile(configFile, globalArguments.ConfigFilePath)

	return true, fmt.Sprintf("=== [region %s] -> [default.region] (%s)", trimmedSelectedRegionResult, globalArguments.CredentialsFilePath)
}
}
