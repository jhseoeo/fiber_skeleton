package middleware

import (
	"github.com/gofiber/fiber/v3"
	fiberRecoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/log"
	"github.com/sirupsen/logrus"
)

func NewRecoverer() fiber.Handler {
	return fiberRecoverer.New(fiberRecoverer.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			log.NewFiberLogEntry(c).WithFields(logrus.Fields{
				"error_stack": e,
			}).Errorln(e)
		},
	})
}
