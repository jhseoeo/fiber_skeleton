package middleware

import (
	"io"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/log"
	"github.com/sirupsen/logrus"
)

func NewLogger() fiber.Handler {
	return logger.New(logger.Config{
		Stream: io.Discard,
		Done: func(c fiber.Ctx, _ []byte) {
			if c.Path() == "/health" {
				// do nothing
				return
			}
			status := c.Response().StatusCode()
			// if an error occurred, do nothing. it will be handled by the error handler
			if 400 <= status && status <= 599 {
				return
			}

			log.NewFiberLogEntry(c).WithFields(logrus.Fields{
				"status": status,
			}).Info("request")

		},
	})
}
