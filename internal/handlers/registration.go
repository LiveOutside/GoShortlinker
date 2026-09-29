package handlers

import (
	"fmt"
	usersdto "goshortlinker/internal/data/dto/users"
	registrationservice "goshortlinker/internal/services/registration"

	"github.com/gofiber/fiber/v3"
)

type RegistrationHandler struct {
	service *registrationservice.Service
}

func NewRegistrationHandler(s *registrationservice.Service) *RegistrationHandler {
	return &RegistrationHandler{
		service: s,
	}
}

func (h *RegistrationHandler) Get(c fiber.Ctx) error {
	return c.Render("register", fiber.Map{})
}

func (h *RegistrationHandler) Post(c fiber.Ctx) error {
	var payload usersdto.RegistrationRequest

	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid parse body")
	}

	user, err := h.service.PostUser(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"redirect_to": fmt.Sprintf("/activate?user_id=%d", user.ID),
	})
}
