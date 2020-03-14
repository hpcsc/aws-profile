package main

import (
	"fmt"
	"github.com/hpcsc/aws-profile/handlers"
	"github.com/hpcsc/aws-profile/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"runtime"
	"strings"
)

func createHandlerMap(app *kingpin.Application, logger utils.Logger) map[string]utils.Handler {
	getHandler := handlers.NewGetHandler(
		app,
		logger,
		utils.GetAWSCallerIdentity,
		utils.ReadCachedCallerIdentity,
		utils.WriteCachedCallerIdentity,
	)
	setHandler := handlers.NewSetHandler(app, utils.SelectProfileFromList, utils.WriteToFile)
	exportHandler := handlers.NewExportHandler(
		app,
		runtime.GOOS == "windows",
		utils.SelectProfileFromList,
		utils.GetAWSCredentials,
	)
	versionHandler := handlers.NewVersionHandler(app)

	return map[string]utils.Handler{
		getHandler.SubCommand.FullCommand():     getHandler,
		setHandler.SubCommand.FullCommand():     setHandler,
		exportHandler.SubCommand.FullCommand():  exportHandler,
		versionHandler.SubCommand.FullCommand(): versionHandler,
	}
}

func main() {
	logger := utils.NewLogrusLogger()

	app := kingpin.New("aws-profile", "simple tool to help switching among AWS profiles more easily")
	app.HelpFlag.Short('h')
	credentialsPathFlag := app.Flag("credentials-path", "Path to AWS Credentials file").Default("~/.aws/credentials").String()
	configPathFlag := app.Flag("config-path", "Path to AWS Config file").Default("~/.aws/config").String()
	handlerMap := createHandlerMap(app, logger)

	if len(os.Args) < 2 {
		app.Usage([]string{})
		os.Exit(1)
	}

	parsedInput := kingpin.MustParse(app.Parse(os.Args[1:]))

	if handler, ok := handlerMap[parsedInput]; ok {
		globalArguments := utils.GlobalArguments{
			CredentialsFilePath: credentialsPathFlag,
			ConfigFilePath:      configPathFlag,
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
