package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/dto/resp"
)

// NewRateLimiter returns a middleware that limits each IP to max requests per
// expiration window. Requests exceeding the limit receive HTTP 429.
func NewRateLimiter(cfg limiter.Config) fiber.Handler {
	cfg.LimitReached = func(c fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(resp.CommonResp{
			Code:    errorcode.ErrTooManyRequests,
			Message: "too many requests",
		})
	}
	return limiter.New(cfg)
}
