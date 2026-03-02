package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// NewCORS returns a CORS middleware configured to allow all origins by default.
// Override AllowOrigins in production to restrict to trusted domains, e.g.:
//
//	middleware.NewCORS(cors.Config{AllowOrigins: []string{"https://example.com"}})
func NewCORS(cfgs ...cors.Config) fiber.Handler {
	if len(cfgs) > 0 {
		return cors.New(cfgs[0])
	}
	return cors.New()
}
