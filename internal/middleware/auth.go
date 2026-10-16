package middleware

import (
	authservice "goshortlinker/internal/services/auth"

	"github.com/gofiber/fiber/v3"
)

const (
	CookieName = "access_token"
	localsKey  = "identity"
)

func Authenticate(tokens *authservice.TokenManager) fiber.Handler {
	return func(c fiber.Ctx) error {
		if raw := c.Cookies(CookieName); raw != "" {
			if identity, err := tokens.Parse(raw); err == nil {
				c.Locals(localsKey, identity)
			}
		}
		return c.Next()
	}
}

func RequireAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		if _, ok := Current(c); ok {
			return c.Next()
		}
		if c.Method() == fiber.MethodGet {
			return c.Redirect().To("/auth")
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":       "требуется авторизация",
			"redirect_to": "/auth",
		})
	}
}

func RedirectIfAuthenticated() fiber.Handler {
	return func(c fiber.Ctx) error {
		if _, ok := Current(c); ok {
			return c.Redirect().To("/")
		}
		return c.Next()
	}
}

func Current(c fiber.Ctx) (authservice.Identity, bool) {
	identity, ok := c.Locals(localsKey).(authservice.Identity)
	return identity, ok
}

func UserID(c fiber.Ctx) (int32, bool) {
	identity, ok := Current(c)
	return identity.UserID, ok
}
