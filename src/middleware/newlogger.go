package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"github.com/jhseoeo/fiber-skeleton/src/pkg/log"
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

		entry := log.NewFiberLogEntry(c).WithFields(logrus.Fields{
			"status":     status,
			"latency_ms": time.Since(start).Milliseconds(),
			"bytes_sent": len(c.Response().Body()),
			"ip":         c.IP(),
		})
		// 2xx → Debug (high-volume, low-value in production)
		// 3xx → Info
		// others logged by the error handler; skip here
		switch {
		case status < 300:
			entry.Debug("request")
		default:
			entry.Info("request")
		}

		return err
	}
}
