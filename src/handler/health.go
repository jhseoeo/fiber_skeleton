package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/dto/resp"
)

// Liveness godoc
//
//	@Summary		Liveness probe
//	@Description	Returns 200 OK when the process is running. Used by orchestrators to
//	@Description	decide whether to restart the container.
//	@Tags			system
//	@Produce		plain
//	@Success		200
//	@Router			/health/live [get]
func Liveness(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

// Readiness godoc
//
//	@Summary		Readiness probe
//	@Description	Returns 200 OK when the service is ready to accept traffic. Add dependency
//	@Description	checks (DB, cache, etc.) here and return 503 when they are unavailable.
//	@Tags			system
//	@Produce		plain
//	@Success		200
//	@Failure		503
//	@Router			/health/ready [get]
func Readiness(c fiber.Ctx) error {
	// TODO: check external dependencies (DB, cache, etc.) and return 503 if unavailable.
	return c.SendStatus(fiber.StatusOK)
}

// NotFound is a catch-all handler that returns a JSON 404 in the CommonResp
// format for any unregistered route.
func NotFound(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(resp.CommonResp{
		Code:    errorcode.ErrNotFound,
		Message: "route not found",
	})
}
