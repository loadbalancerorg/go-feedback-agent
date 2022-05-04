package main

import (
	log "github.com/sirupsen/logrus"
)

type Logging struct {
	Logger *log.Logger
}

func setupLogging() {
	newLogger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	eventLog = newLogger
}

func (l *Logging) SetLogLevel(level string) {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		panic(err.Error())
	}
	eventLog.Logger.SetLevel(logLevel)
	if logLevel == log.DebugLevel {
		eventLog.Logger.Infoln("**** DEBUG LOGGING ENABLED ****")
	}
}
