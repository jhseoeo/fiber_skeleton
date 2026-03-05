package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/typeerr"
)

// NewTimeout returns a middleware that cancels the request context after d.
// All downstream calls that respect context cancellation (DB queries, HTTP
// clients, etc.) will be aborted automatically when the deadline is exceeded.
//
// If d <= 0 the middleware is a no-op.
func NewTimeout(d time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		if d <= 0 {
			return c.Next()
		}

		ctx, cancel := context.WithTimeout(c.Context(), d)
		defer cancel()

		c.SetContext(ctx)

		err := c.Next()
		if ctx.Err() == context.DeadlineExceeded {
			return typeerr.NewErrorResp(ctx.Err(), errorcode.ErrRequestTimeout, "request timeout")
		}
		return err
	}
}
