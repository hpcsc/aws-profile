package handlers

import (
	"errors"
	"fmt"
	"github.com/hpcsc/aws-profile/internal/utils"
	"strings"

	"github.com/hpcsc/aws-profile/internal/awsconfig"
	"github.com/hpcsc/aws-profile/internal/io"
	"gopkg.in/alecthomas/kingpin.v2"
)

type SelectRegionFn func([]string, string) ([]byte, error)

type SetRegionHandler struct {
	SubCommand   *kingpin.CmdClause
	SelectRegion SelectRegionFn
	WriteToFile  WriteToFileFn
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

	selectRegionResult, err := handler.SelectRegion(getAllRegions(), "Select an AWS region")
	var cancelled *utils.CancelledError
	if errors.As(err, &cancelled) {
		return true, ""
	}

	if err != nil {
		return false, fmt.Sprintf("Failed to select region: %v", err)
	}

	trimmedSelectedRegionResult := strings.TrimSuffix(string(selectRegionResult), "\n")

	awsconfig.SetSelectedRegionAsDefault(trimmedSelectedRegionResult, configFile)
	if err := handler.WriteToFile(configFile, globalArguments.ConfigFilePath); err != nil {
		return false, err.Error()
	}

	return true, fmt.Sprintf("=== [region %s] -> [default.region] (%s)", trimmedSelectedRegionResult, globalArguments.ConfigFilePath)
}

func getAllRegions() []string {
	return []string{
		"af-south-1",
		"ap-east-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"cn-north-1",
		"cn-northwest-1",
		"eu-central-1",
		"eu-north-1",
		"eu-south-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"me-south-1",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-gov-east-1",
		"us-gov-west-1",
		"us-west-1",
		"us-west-2",
	}
}
