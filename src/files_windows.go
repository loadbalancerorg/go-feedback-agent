//go:build windows
// +build windows

package main

import (
	"github.com/Freman/eventloghook"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc/eventlog"
	"os"
)

const CONFIG_FILE = "C:/ProgramData/LoadBalancer.org/LoadBalancer/config.xml"
const LOG_FILE = "C:/ProgramData/LoadBalancer.org/LoadBalancer/lbfbalogfile"

// NewLogger Returns a new distributes logging object
func NewLogger() (Logging, error) {
	logging := Logging{
		log.New(),
	}
	// Setup event log
	const name = "Feedback-Agent"
	const supports = eventlog.Error | eventlog.Warning | eventlog.Info
	err := eventlog.InstallAsEventCreate(name, supports)
	if err != nil {
		return Logging{}, err
	}
	defer func() {
		err = eventlog.Remove(name)
		if err != nil {
			return
		}
	}()
	eventLog, err := eventlog.Open(name)
	if err != nil {
		return Logging{}, err
	}
	defer eventLog.Close()
	hook := eventloghook.NewHook(eventLog)
	// Attach our event log hook
	logging.Logger.Hooks.Add(hook)

	f, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	logging.Logger.SetOutput(f)
	//logging.Logger.SetOutput(io.Discard)

	// Return with our logger
	return logging, nil
}

func localCMD() (string, string) {
	return "cmd", "/c"
}
