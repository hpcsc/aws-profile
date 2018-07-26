package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/hpcsc/aws-profile-utils/handlers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("sub command is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		getCommand := handlers.NewGetHandler()
		getCommand.Handle(os.Args[2:])
	case "set":
		setCommand := handlers.NewSetHandler()
		setCommand.Handle(os.Args[2:])
	case "version":
		versionCommand := handlers.NewVersionHandler()
		versionCommand.Handle(os.Args[2:])
	default:
		flag.Usage()
		os.Exit(1)
	}
}