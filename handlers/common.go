package handlers

import (
	"os/user"
	"strings"
	"path/filepath"
	"gopkg.in/ini.v1"
		)

func ReadFile(filePath string) (*ini.File, error) {
	path := ExpandHomeDirectory(filePath)
	return ini.Load(path)
}

func ExpandHomeDirectory(filePath string) string {
	usr, _ := user.Current()
	homeDirectory := usr.HomeDir
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(homeDirectory, filePath[2:])
	}

	return filePath
}

type Handler interface {
	Handle()
}