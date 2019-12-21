package main

import (
	"fmt"
	"github.com/hpcsc/aws-profile-utils/handlers"
	"github.com/hpcsc/aws-profile-utils/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

func createHandlerMap(app *kingpin.Application) map[string]utils.Handler {
	getHandler := handlers.NewGetHandler(app)
	setHandler := handlers.NewSetHandler(app, utils.SelectProfileFromList, utils.WriteToFile)
	versionHandler := handlers.NewVersionHandler(app)

	return map[string]utils.Handler{
		getHandler.SubCommand.FullCommand(): getHandler,
		setHandler.SubCommand.FullCommand(): setHandler,
		versionHandler.SubCommand.FullCommand(): versionHandler,
	}
}

func main() {
	app := kingpin.New("aws-profile-utils", "simple tool to help switching among AWS profiles more easily")
	app.HelpFlag.Short('h')
	handlerMap := createHandlerMap(app)

	if len(os.Args) < 2 {
		app.Usage([]string {})
		os.Exit(1)
	}

	parsedInput := kingpin.MustParse(app.Parse(os.Args[1:]))

	if handler, ok := handlerMap[parsedInput]; ok {
		success, message := handler.Handle()
		if !strings.EqualFold(message, "") {
			fmt.Println(message)
		}

		if success {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		app.Usage([]string {})
	}
}
