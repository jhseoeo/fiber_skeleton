package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/dto/resp"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/log"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/typeerr"
	"github.com/sirupsen/logrus"
)

func NewErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		var fiberError *fiber.Error
		if errors.As(err, &fiberError) {
			code := errorcode.ErrorCode(fiberError.Code * 100)
			return c.Status(fiberError.Code).JSON(resp.CommonResp{
				Code:    code,
				Message: fiberError.Message,
			})
		}

		var (
			e      typeerr.ErrorResp
			status = fiber.StatusInternalServerError
			code   = errorcode.ErrInternalServer
			inner  error
		)
		if errors.As(err, &e) {
			status = e.Code.HTTPStatus()
			code = e.Code
			inner = e.Err
		} else {
			inner = err
		}

		entry := log.NewFiberLogEntry(c).WithFields(logrus.Fields{
			"status":  status,
			"message": e.Message,
		}).WithError(inner)
		if status >= 500 {
			entry.Errorln(e.Message)
		} else {
			entry.Infoln(e.Message)
		}

		return c.Status(status).JSON(resp.CommonResp{
			Code:    code,
			Message: e.Message,
		})
	}
}
