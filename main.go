package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/hpcsc/aws-profile-utils/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("sub command is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		getCommand := commands.NewGetCommand()
		getCommand.Handle(os.Args[2:])
	case "set":
		setCommand := commands.NewSetCommand()
		setCommand.Handle(os.Args[2:])
	case "version":
		versionCommand := commands.NewVersionCommand()
		versionCommand.Handle(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}