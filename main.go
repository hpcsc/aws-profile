package main

import (
		"fmt"
	"os"
	"github.com/hpcsc/aws-profile-utils/handlers"
		"strings"
)

var handlerMap = map[string]handlers.Handler {
	"get": handlers.NewGetHandler(),
	"set": handlers.NewSetHandler(),
	"version": handlers.NewVersionHandler(),
}

func printUsage() {
	var commandNames []string

	for name := range handlerMap {
		commandNames = append(commandNames, name)
	}

	fmt.Printf("Usage: %s [%s] arguments... \n", os.Args[0], strings.Join(commandNames[:], "|"))
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if handler, ok := handlerMap[os.Args[1]]; ok {
		handler.Handle(os.Args[2:])
	} else {
		printUsage()
		os.Exit(1)
	}
}