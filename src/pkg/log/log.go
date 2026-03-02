package log

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/config"
	"github.com/sirupsen/logrus"
)

func Init(cfg *config.Config) {
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	if cfg.Env == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}
}

func NewFiberLogEntry(c fiber.Ctx) *logrus.Entry {
	fields := logrus.Fields{
		"method": c.Method(),
		"path":   c.Path(),
	}
	if requestID := c.Locals("requestid"); requestID != nil {
		fields["request_id"] = requestID
	}
	return logrus.WithFields(fields)
}
