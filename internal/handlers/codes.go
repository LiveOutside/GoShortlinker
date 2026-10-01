package handlers

import (
	"errors"
	dtocodes "goshortlinker/internal/data/dto/activation_codes"
	activationservice "goshortlinker/internal/services/mailer"

	"github.com/gofiber/fiber/v3"
)

type MailerHandler struct {
	service *activationservice.Service
}

func NewMailerHandler(s *activationservice.Service) *MailerHandler {
	return &MailerHandler{
		service: s,
	}
}

func (h *MailerHandler) Get(c fiber.Ctx) error {
	return c.Render("activate", fiber.Map{})
}

func (h *MailerHandler) Post(c fiber.Ctx) error {
	var payload dtocodes.ActivateRequest
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtocodes.ActivateResponse{
			Error: "invalid body",
		})
	}

	err := h.service.Verify(c.Context(), payload)
	if err != nil {
		switch err {
		case activationservice.ErrCodeExpired:
			return c.Status(fiber.StatusBadRequest).JSON(dtocodes.ActivateResponse{
				Expired: true,
				Error:   "code expired, new code sent",
			})
		case activationservice.ErrInvalidCode:
			return c.Status(fiber.StatusBadRequest).JSON(dtocodes.ActivateResponse{
				Error: "invalid code",
			})
		case activationservice.ErrCodeNotFound:
			return c.Status(fiber.StatusNotFound).JSON(dtocodes.ActivateResponse{
				Error: "no activation pending",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(dtocodes.ActivateResponse{
				Error: "internal error",
			})
		}
	}

	if err := h.service.Delete(c.Context(), payload.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtocodes.ActivateResponse{
			Error: "internal error",
		})
	}

	return c.JSON(dtocodes.ActivateResponse{RedirectTo: "/auth"})
}

func (h *MailerHandler) Resend(c fiber.Ctx) error {
	var payload dtocodes.ResendRequest

	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtocodes.ActivateResponse{
			Error: "invalid body",
		})
	}

	if err := h.service.Resend(c.Context(), payload.UserID); err != nil {
		if errors.Is(err, activationservice.ErrCodeNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dtocodes.ActivateResponse{
				Error: "no activation pending",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dtocodes.ActivateResponse{
			Error: "internal error",
		})
	}

	return c.JSON(dtocodes.ActivateResponse{})
}
