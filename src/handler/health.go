package handler

import "github.com/gofiber/fiber/v3"

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Returns 200 OK when the server is healthy
//	@Tags			system
//	@Produce		plain
//	@Success		200
//	@Router			/health [get]
func HealthCheck(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
