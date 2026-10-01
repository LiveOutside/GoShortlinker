package handlers

import "github.com/gofiber/fiber/v3"

type AuthenticationHandler struct {
}

func NewAuthenticationHandler() *AuthenticationHandler {
	return &AuthenticationHandler{}
}

func (h *AuthenticationHandler) Get(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}
