package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/log"
	"github.com/sirupsen/logrus"
)

func NewLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		if c.Path() == "/health/live" || c.Path() == "/health/ready" {
			return err
		}

		status := c.Response().StatusCode()
		// if an error occurred, do nothing. it will be handled by the error handler
		if 400 <= status && status <= 599 {
			return err
		}

		log.NewFiberLogEntry(c).WithFields(logrus.Fields{
			"status":     status,
			"latency_ms": time.Since(start).Milliseconds(),
			"bytes_sent": len(c.Response().Body()),
			"ip":         c.IP(),
		}).Info("request")

		return err
	}
}
