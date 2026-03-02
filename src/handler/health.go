package handler

import "github.com/gofiber/fiber/v3"

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
