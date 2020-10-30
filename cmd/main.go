package main

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/aws"
	"github.com/hpcsc/aws-profile/internal/handlers"
	"github.com/hpcsc/aws-profile/internal/io"
	"github.com/hpcsc/aws-profile/internal/log"
	"github.com/hpcsc/aws-profile/internal/tui"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"runtime"
	"strings"
)

func createHandlerMap(app *kingpin.Application, logger log.Logger) map[string]handlers.Handler {
	isWindows := runtime.GOOS == "windows"

	getHandler := handlers.NewGetHandler(
		app,
		logger,
		aws.GetAWSCallerIdentity,
		io.ReadCachedCallerIdentity,
		io.WriteCachedCallerIdentity,
	)
	setHandler := handlers.NewSetHandler(app, tui.SelectProfileFromList, io.WriteToFile)
	exportHandler := handlers.NewExportHandler(
		app,
		isWindows,
		tui.SelectProfileFromList,
		aws.GetAWSCredentials,
	)
	unsetHandler := handlers.NewUnsetHandler(app, isWindows)
	versionHandler := handlers.NewVersionHandler(app)

	return map[string]handlers.Handler{
		getHandler.SubCommand.FullCommand():     getHandler,
		setHandler.SubCommand.FullCommand():     setHandler,
		exportHandler.SubCommand.FullCommand():  exportHandler,
		unsetHandler.SubCommand.FullCommand():   unsetHandler,
		versionHandler.SubCommand.FullCommand(): versionHandler,
	}
}

func main() {
	logger := log.NewLogrusLogger()

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
		globalArguments := handlers.GlobalArguments{
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
