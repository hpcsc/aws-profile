package main

import (
			"os"
	"github.com/hpcsc/aws-profile-utils/handlers"
			"gopkg.in/alecthomas/kingpin.v2"
	)

func createHandlerMap(app *kingpin.Application) map[string]handlers.Handler{
	getHandler := handlers.NewGetHandler(app)
	setHandler := handlers.NewSetHandler(app)
	versionHandler := handlers.NewVersionHandler(app)

	return map[string]handlers.Handler {
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
		handler.Handle()
	} else {
		app.Usage([]string {})
	}
}