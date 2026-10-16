package handlers

import (
	"errors"
	"fmt"
	dtousers "goshortlinker/internal/data/dto/users"
	"goshortlinker/internal/middleware"
	authservice "goshortlinker/internal/services/auth"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

var validate = validator.New()

type AuthenticationHandler struct {
	service      *authservice.Service
	secureCookie bool
}

func NewAuthenticationHandler(s *authservice.Service) *AuthenticationHandler {
	return &AuthenticationHandler{
		service:      s,
		secureCookie: os.Getenv("COOKIE_SECURE") != "false",
	}
}

func (h *AuthenticationHandler) Get(c fiber.Ctx) error {
	if _, ok := middleware.UserID(c); ok {
		return c.Redirect().To("/")
	}
	return c.Render("login", fiber.Map{})
}

func (h *AuthenticationHandler) Post(c fiber.Ctx) error {
	var payload dtousers.AuthRequest

	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "некорректный запрос"})
	}
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "введите корректные email и пароль"})
	}

	res, err := h.service.Login(c.Context(), payload)
	switch {
	case errors.Is(err, authservice.ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "неверный email или пароль"})
	case errors.Is(err, authservice.ErrNotActivated):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":       "аккаунт не активирован",
			"redirect_to": fmt.Sprintf("/activate?user_id=%d", res.UserID),
		})
	case err != nil:
		log.Errorf("login failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "внутренняя ошибка сервера"})
	}

	h.setAuthCookie(c, res.Token, res.ExpiresAt)
	return c.JSON(fiber.Map{"redirect_to": "/"})
}

func (h *AuthenticationHandler) Logout(c fiber.Ctx) error {
	h.clearAuthCookie(c)
	return c.JSON(fiber.Map{"redirect_to": "/"})
}

func (h *AuthenticationHandler) setAuthCookie(c fiber.Ctx, token string, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     middleware.CookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
		HTTPOnly: true,
		Secure:   h.secureCookie,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func (h *AuthenticationHandler) clearAuthCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     middleware.CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.secureCookie,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}
