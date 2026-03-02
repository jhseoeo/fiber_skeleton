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
			entry := log.NewFiberLogEntry(c).WithFields(logrus.Fields{
				"status":  fiberError.Code,
				"message": fiberError.Message,
			})
			if fiberError.Code >= 500 {
				entry.Errorln(fiberError.Message)
			} else {
				entry.Infoln(fiberError.Message)
			}
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
			msg    string
			data   any
		)
		if errors.As(err, &e) {
			status = e.Code.HTTPStatus()
			code = e.Code
			inner = e.Err
			msg = e.Message
			data = e.Data
		} else {
			inner = err
			msg = "internal server error"
		}

		entry := log.NewFiberLogEntry(c).WithFields(logrus.Fields{
			"status":  status,
			"message": msg,
		}).WithError(inner)
		if status >= 500 {
			entry.Errorln(msg)
		} else {
			entry.Infoln(msg)
		}

		return c.Status(status).JSON(resp.CommonResp{
			Code:    code,
			Message: msg,
			Data:    data,
		})
	}
}
