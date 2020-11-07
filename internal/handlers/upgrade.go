package handlers

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/log"
	"github.com/hpcsc/aws-profile/internal/upgrade"
	"github.com/hpcsc/aws-profile/internal/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type UpgradeHandler struct {
	SubCommand *kingpin.CmdClause
	Logger     log.Logger
	Arguments  UpgradeCommandArguments
}

type UpgradeCommandArguments struct {
	IncludePrerelease *bool
}

func NewUpgradeHandler(app *kingpin.Application, logger log.Logger, ) UpgradeHandler {
	subCommand := app.Command("upgrade", "upgrade to latest stable version")

	includePrerelease := subCommand.Flag("prerelease", "Include prerelease").Default("false").Bool()

	return UpgradeHandler{
		SubCommand: subCommand,
		Arguments: UpgradeCommandArguments{
			IncludePrerelease: includePrerelease,
		},
		Logger: logger,
	}
}

func (handler UpgradeHandler) Handle(globalArguments GlobalArguments) (bool, string) {
	binaryPath, err := os.Executable()
	if err != nil {
		return false, fmt.Sprintf("failed to get current executable path: %v", err)
	}

	message, err := upgrade.ToLatest(binaryPath, *handler.Arguments.IncludePrerelease, version.Current())
	if err != nil {
		return false, err.Error()
	}

	return true, message
}
