package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// corsMaxAge is the default preflight cache duration in seconds (1 hour).
const corsMaxAge = 3600

// NewCORS returns a CORS middleware configured to allow all origins by default.
// Override AllowOrigins in production to restrict to trusted domains, e.g.:
//
//	middleware.NewCORS(cors.Config{AllowOrigins: []string{"https://example.com"}})
//
// MaxAge defaults to corsMaxAge seconds (1 hour) to reduce preflight requests.
func NewCORS(cfgs ...cors.Config) fiber.Handler {
	cfg := cors.Config{MaxAge: corsMaxAge}
	if len(cfgs) > 0 {
		cfg = cfgs[0]
		if cfg.MaxAge == 0 {
			cfg.MaxAge = corsMaxAge
		}
	}
	return cors.New(cfg)
}
