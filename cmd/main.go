package main

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/config"
	"os"
	"runtime"
	"strings"

	"github.com/hpcsc/aws-profile/internal/aws"
	"github.com/hpcsc/aws-profile/internal/handlers"
	"github.com/hpcsc/aws-profile/internal/io"
	"github.com/hpcsc/aws-profile/internal/log"
	"github.com/hpcsc/aws-profile/internal/tui"
	"github.com/hpcsc/aws-profile/internal/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

func createHandlerMap(app *kingpin.Application, logger log.Logger, config *config.Config) map[string]handlers.Handler {
	isWindows := runtime.GOOS == "windows"

	getHandler := handlers.NewGetHandler(
		app,
		logger,
		aws.GetAWSCallerIdentity,
		io.ReadCachedCallerIdentity,
		io.WriteCachedCallerIdentity,
	)
	setHandler := handlers.NewSetHandler(app, tui.SelectProfileFromList, io.WriteToFile)
	setRegionHandler := handlers.NewSetRegionHandler(app, config, tui.SelectValueFromList, io.WriteToFile)
	exportHandler := handlers.NewExportHandler(
		app,
		isWindows,
		tui.SelectProfileFromList,
		aws.GetAWSCredentials,
	)
	unsetHandler := handlers.NewUnsetHandler(app, isWindows)
	upgradeHandler := handlers.NewUpgradeHandler(app, logger)
	versionHandler := handlers.NewVersionHandler(app)

	return map[string]handlers.Handler{
		getHandler.SubCommand.FullCommand():       getHandler,
		setHandler.SubCommand.FullCommand():       setHandler,
		setRegionHandler.SubCommand.FullCommand(): setRegionHandler,
		exportHandler.SubCommand.FullCommand():    exportHandler,
		unsetHandler.SubCommand.FullCommand():     unsetHandler,
		upgradeHandler.SubCommand.FullCommand():   upgradeHandler,
		versionHandler.SubCommand.FullCommand():   versionHandler,
	}
}

func main() {
	logger := log.NewLogrusLogger()

	configPath := utils.GetEnvVariableOrDefault("AWS_PROFILE_CONFIG", "~/.aws-profile/config.yaml")
	config, err := config.FromFile(utils.ExpandHomeDirectory(configPath))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	app := kingpin.New("aws-profile", "simple tool to help switching among AWS profiles more easily")
	app.HelpFlag.Short('h')
	handlerMap := createHandlerMap(app, logger, config)

	if len(os.Args) < 2 {
		app.Usage([]string{})
		os.Exit(1)
	}

	parsedInput := kingpin.MustParse(app.Parse(os.Args[1:]))

	if handler, ok := handlerMap[parsedInput]; ok {
		globalArguments := handlers.GlobalArguments{
			CredentialsFilePath: utils.GetEnvVariableOrDefault("AWS_SHARED_CREDENTIALS_FILE", "~/.aws/credentials"),
			ConfigFilePath:      utils.GetEnvVariableOrDefault("AWS_CONFIG_FILE", "~/.aws/config"),
		}

		success, message := handler.Handle(globalArguments)
		if !strings.EqualFold(message, "") {
			fmt.Println(message)
		}

		if success {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		app.Usage([]string{})
	}
}
