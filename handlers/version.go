package handlers

import (
	"flag"
	"fmt"
		)

type VersionHandler struct {
	FlagSet *flag.FlagSet
}

func NewVersionHandler() VersionHandler {
	flagSet := flag.NewFlagSet("version", flag.ExitOnError)

	return VersionHandler{
		FlagSet: flagSet,
	}
}

var version = "undefined"
func (handler VersionHandler) Handle(arguments []string) {
	fmt.Printf("aws-profile-utils (v%s)", version)
}
