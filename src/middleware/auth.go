package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/typeerr"
)

// Claims holds the JWT registered claims plus any custom application claims.
// TODO: add custom fields such as UserID, Role, etc.
type Claims struct {
	jwt.RegisteredClaims
	// UserID uint   `json:"user_id"`
	// Role   string `json:"role"`
}

const claimsKey = "claims"

// NewAuthMiddleware returns a JWT authentication middleware.
// It validates the Bearer token from the Authorization header and
// stores the parsed Claims in the request context via c.Locals.
//
// TODO: load jwtSecret from config and pass it in.
func NewAuthMiddleware(jwtSecret []byte) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return typeerr.NewErrorResp(
				fmt.Errorf("missing authorization header"),
				errorcode.ErrUnauthorized,
				"missing or invalid authorization header",
			)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return typeerr.NewErrorResp(err, errorcode.ErrUnauthorized, "invalid or expired token")
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return typeerr.NewErrorResp(fmt.Errorf("invalid claims type"), errorcode.ErrUnauthorized, "invalid token claims")
		}

		c.Locals(claimsKey, claims)
		return c.Next()
	}
}

// GetClaims extracts the parsed JWT claims stored by NewAuthMiddleware.
// Returns (nil, false) if the middleware was not applied or the token was invalid.
func GetClaims(c fiber.Ctx) (*Claims, bool) {
	claims, ok := c.Locals(claimsKey).(*Claims)
	return claims, ok
}
