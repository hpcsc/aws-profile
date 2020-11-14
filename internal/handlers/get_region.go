package handlers

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/io"
	"gopkg.in/alecthomas/kingpin.v2"
)

type GetRegionHandler struct {
	SubCommand *kingpin.CmdClause
}

func NewGetRegionHandler(app *kingpin.Application) GetRegionHandler {
	subCommand := app.Command("get-region", "get current region set in default profile")

	return GetRegionHandler{
		SubCommand: subCommand,
	}
}

func (handler GetRegionHandler) Handle(globalArguments GlobalArguments) (bool, string) {
	configFile, err := io.ReadFile(globalArguments.ConfigFilePath)
	if err != nil {
		return false, fmt.Sprintf("Fail to read AWS config file: %v", err)
	}

	defaultProfileInConfig := configFile.Section("default")
	if defaultProfileInConfig.HasKey("region") && defaultProfileInConfig.Key("region").Value() != "" {
		return true, defaultProfileInConfig.Key("region").Value()
	}

	return true, "no region set"
}
