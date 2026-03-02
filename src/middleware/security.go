package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

// NewSecurity returns a middleware that sets common security headers using
// Fiber's built-in helmet middleware.
//
// Headers set by default:
//   - X-XSS-Protection: 0
//   - X-Content-Type-Options: nosniff
//   - X-Frame-Options: SAMEORIGIN
//   - Content-Security-Policy: default-src 'self'
//   - Strict-Transport-Security: max-age=31536000; includeSubDomains
func NewSecurity() fiber.Handler {
	return helmet.New()
}
