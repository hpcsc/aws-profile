package handlers

import "gopkg.in/ini.v1"

type WriteToFileFn func(*ini.File, string) error
