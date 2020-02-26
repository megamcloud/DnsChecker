package main

import (
	"os"

	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func newLogger(config *config) *logrus.Entry {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Debug output if app is running in debug mode
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if config.LogglyToken != "" {
		loggly := logrusly.NewLogglyHook(config.LogglyToken, config.ApplicationName, log.InfoLevel, config.ApplicationName)
		log.AddHook(loggly)
	} else {
		log.Info("No Loggly token found, only logging to console")
	}

	return log.WithField("Application", config.ApplicationName)
}
