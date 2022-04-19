//go:build windows
// +build windows

package main

import (
	"github.com/Freman/eventloghook"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc/eventlog"
)

const CONFIG_FILE = "C:/ProgramData/LoadBalancer.org/LoadBalancer/config.xml"

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
	// Setup event log
	eventLog, err := eventlog.Open("Feedback Agent")
	if err != nil {
		return Logging{}, err
	}
	defer eventLog.Close()
	hook := eventloghook.NewHook(eventLog)
	// Attach our event log hook
	logging.Logger.Hooks.Add(hook)
	// Return with our logger
	return logging, nil
}
