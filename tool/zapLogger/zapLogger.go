// Package zapLogger provide pre-configured zap loggers.
package zapLogger

import (
	"log"

	"encoding/json"

	"go.uber.org/zap"
)

// NewServerLogger create zap config specific for GoK service.
func NewServerLogger(debug bool) *zap.Logger {
	level := "warn"
	if debug {
		level = "debug"
	}

	rawJSON := []byte(`{
	  "level": "` + level + `",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/gok-service.log"],
	  "errorOutputPaths": ["stderr", "/tmp/gok-service.error.log"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		log.Fatal("Failed to unmarshal logger config: %w", err)
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatal("Failed to init logger: %w", err)
	}
	logger.Debug("Service logger created successfully")

	return logger
}
