package main

import (
	log "github.com/sirupsen/logrus"
)

type Logging struct {
	Logger *log.Logger
	Level  LogsLevel
}

// Debug Send a debug message to the cluster
func (l *Logging) Debug(message string) {
	if l.Level == LogsDEBUG {
		l.Logger.Debugf("%s", message)
	}
}

// Warn Send a warning message to the cluster
func (l *Logging) Warn(message string) {
	l.Logger.Warnf("%s", message)
}

// Info Send an info message to the cluster
func (l *Logging) Info(message string) {
	l.Logger.Infof("%s", message)
}

// Info Send an error message to the cluster
func (l *Logging) Error(message string) {
	l.Logger.Errorf("%s", message)
}
