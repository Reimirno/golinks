package logging

import (
	"go.uber.org/zap"
)

var loggerSet *zap.Logger

func Initialize(debug bool) error {
	var err error
	var logger *zap.Logger
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	loggerSet = logger
	defer loggerSet.Sync()
	return nil
}

func NewLogger(name string) *zap.SugaredLogger {
	if loggerSet == nil {
		Initialize(true)
	}
	return loggerSet.Named(name).Sugar()
}
