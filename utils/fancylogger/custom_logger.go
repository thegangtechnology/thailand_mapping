package fancylogger

import (
	"go-template/config"
	"os"

	log "github.com/sirupsen/logrus"

	joonix "github.com/joonix/log"
	"github.com/labstack/echo/v4"
	echologrus "github.com/spirosoik/echo-logrus"
)

func SetupLogger(e *echo.Echo, gcp bool) {
	logger := log.New()

	level := os.Getenv(config.LogLevel)

	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
		logger.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
		logger.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
		logger.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
		logger.SetLevel(log.InfoLevel)
	}

	if gcp {
		gcpSafeFormatter := joonix.NewFormatter()

		logger.SetFormatter(gcpSafeFormatter)
		log.SetFormatter(gcpSafeFormatter)
	} else {
		customFormatter := new(log.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02 15:04:05"
		customFormatter.FullTimestamp = true

		logger.SetFormatter(customFormatter)
		log.SetFormatter(customFormatter)
	}

	if e != nil {
		echoLogger := echologrus.NewLoggerMiddleware(logger)
		e.Logger = echoLogger
		e.Use(echoLogger.Hook())
	}
}
