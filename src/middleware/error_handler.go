package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/log"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/typeerr"
	"github.com/sirupsen/logrus"
)

func NewErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		var fiberError *fiber.Error
		if errors.As(err, &fiberError) {
			return c.Status(fiberError.Code).JSON(fiber.Map{"message": fiberError.Message})
		}

		var (
			e       typeerr.ErrorResp
			message string
			status  = fiber.StatusInternalServerError
			inner   error
		)
		if errors.As(err, &e) {
			message = e.Message
			status = e.Status
			inner = e.Err
		} else {
			message = err.Error()
			inner = err
		}

		entry := log.NewFiberLogEntry(c).WithFields(logrus.Fields{
			"status":  status,
			"message": message,
		}).WithError(inner)
		if status >= 500 {
			entry.Errorln(message)
		} else {
			entry.Infoln(message)
		}

		return c.Status(status).JSON(fiber.Map{"message": message})
	}
}
