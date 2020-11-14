package handlers

import (
	"errors"
	"fmt"
	"github.com/hpcsc/aws-profile/internal/config"
	"github.com/hpcsc/aws-profile/internal/utils"
	"strings"

	"github.com/hpcsc/aws-profile/internal/awsconfig"
	"github.com/hpcsc/aws-profile/internal/io"
	"gopkg.in/alecthomas/kingpin.v2"
)

type SelectRegionFn func([]string, string, *config.Config) ([]byte, error)

type SetRegionHandler struct {
	SubCommand   *kingpin.CmdClause
	SelectRegion SelectRegionFn
	WriteToFile  WriteToFileFn
	Config       *config.Config
}

func NewSetRegionHandler(app *kingpin.Application, config *config.Config, selectRegionFn SelectRegionFn, writeToFileFn WriteToFileFn) SetRegionHandler {
	subCommand := app.Command("set-region", "set the region of the default profile")

	return SetRegionHandler{
		SubCommand:   subCommand,
		SelectRegion: selectRegionFn,
		WriteToFile:  writeToFileFn,
		Config:       config,
	}
}

func (handler SetRegionHandler) Handle(globalArguments GlobalArguments) (bool, string) {
	configFile, err := io.ReadFile(globalArguments.ConfigFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", err)
	}

	selectRegionResult, err := handler.SelectRegion(handler.Config.Regions, "Select an AWS region", handler.Config)
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
