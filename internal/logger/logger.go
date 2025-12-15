package logger

import (
	"os"
	"path/filepath"

	"github.com/keerthana/linkedin-automation-poc/pkg/models"
	"github.com/sirupsen/logrus"
)

// New creates a new configured logger
func New(config *models.LoggingConfig) (*logrus.Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output
	if config.OutputPath != "" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(config.OutputPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}

		// Open log file
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		logger.SetOutput(file)
	} else {
		logger.SetOutput(os.Stdout)
	}

	return logger, nil
}
