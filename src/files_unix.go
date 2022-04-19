//go:build !windows
// +build !windows

package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ssgreg/journalhook"
)

const CONFIG_FILE = "/usr/local/etc/lbfbaconfig.xml"

func setupLogging() {
	newLogger, err := NewLogger()
	if err != nil {
		panic("unable to create logger")
	}
	eventLog = newLogger
}

// NewLogger Returns a new distributes logging object
func NewLogger() (Logging, error) {
	logging := Logging{
		log.New(),
		LogsINFO,
	}
	// Setup journal
	hook, err := journalhook.NewJournalHook()
	if err != nil {
		return Logging{}, err
	}
	// Attach our journal hook
	logging.Logger.Hooks.Add(hook)
	// Return with our logger
	return logging, nil
}
