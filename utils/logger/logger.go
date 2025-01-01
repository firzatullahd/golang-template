package logger

import (
	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func NewLogger(env string) {
	log := logrus.New()

	var level logrus.Level
	switch env {
	case "production":
		level = logrus.InfoLevel
	default:
		level = logrus.DebugLevel
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus.JSONFormatter{})

	Log = log
}
