package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// NewCORS returns a CORS middleware configured to allow all origins by default.
// Override AllowOrigins in production to restrict to trusted domains, e.g.:
//
//	middleware.NewCORS(cors.Config{AllowOrigins: []string{"https://example.com"}})
//
// MaxAge defaults to 3600 seconds (1 hour) to reduce preflight requests.
func NewCORS(cfgs ...cors.Config) fiber.Handler {
	cfg := cors.Config{MaxAge: 3600}
	if len(cfgs) > 0 {
		cfg = cfgs[0]
		if cfg.MaxAge == 0 {
			cfg.MaxAge = 3600
		}
	}
	return cors.New(cfg)
}
