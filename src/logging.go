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

// Warn Logs a warning message
func (l *Logging) Warn(message string) {
	l.Logger.Warnf("%s", message)
}

// Info Logs an info message
func (l *Logging) Info(message string) {
	l.Logger.Infof("%s", message)
}

// Error Logs an error message
func (l *Logging) Error(message string) {
	l.Logger.Errorf("%s", message)
}

// ErrorErr Logs an error message
func (l *Logging) ErrorErr(err error) {
	l.Logger.Println(err)
}
