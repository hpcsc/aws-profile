package handlers

import (
	"os/user"
	"strings"
	"path/filepath"
	)

func ExpandHomeDirectory(filePath string) string {
	usr, _ := user.Current()
	homeDirectory := usr.HomeDir
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(homeDirectory, filePath[2:])
	}

	return filePath
}

type Handler interface {
	Handle(arguments []string)
}