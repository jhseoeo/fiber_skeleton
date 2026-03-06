package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/typeerr"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/validate"
)

func parseID(c fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	return uint(id), err
}

// bindJSON binds a JSON request body into dst and validates it.
// Returns a ready-to-return ErrorResp on failure, or nil on success.
func bindJSON(c fiber.Ctx, dst any) error {
	if err := c.Bind().JSON(dst); err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidBody, "invalid request body")
	}
	if err := validate.Struct(dst); err != nil {
		return typeerr.NewErrorRespWithData(err, errorcode.ErrInvalidBody, "validation failed", err)
	}
	return nil
}

// bindQuery binds query parameters into dst and validates them.
// Returns a ready-to-return ErrorResp on failure, or nil on success.
func bindQuery(c fiber.Ctx, dst any) error {
	if err := c.Bind().Query(dst); err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrBadRequest, "invalid query parameters")
	}
	if err := validate.Struct(dst); err != nil {
		return typeerr.NewErrorRespWithData(err, errorcode.ErrBadRequest, "validation failed", err)
	}
	return nil
}
