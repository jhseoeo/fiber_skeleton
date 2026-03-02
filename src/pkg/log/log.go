package log

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

func NewFiberLogEntry(c fiber.Ctx) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"method": c.Method(),
		"path":   c.Path(),
	})
}
