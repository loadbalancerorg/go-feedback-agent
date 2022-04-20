//go:build !windows
// +build !windows

package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ssgreg/journalhook"
	"io/ioutil"
)

const CONFIG_FILE = "/usr/local/etc/lbfbaconfig.xml"

// NewLogger Returns a new distributes logging object
func NewLogger() (Logging, error) {
	logging := Logging{
		log.New(),
	}
	// Setup journal
	hook, err := journalhook.NewJournalHook()
	if err != nil {
		return Logging{}, err
	}
	// Attach our journal hook
	logging.Logger.Hooks.Add(hook)

	logging.Logger.SetOutput(ioutil.Discard)

	// Return with our logger
	return logging, nil
}

func localCMD() (string, string) {
	return "/bin/sh", "-c"
}
