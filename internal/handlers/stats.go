package handlers

import (
	"goshortlinker/internal/middleware"
	linkservice "goshortlinker/internal/services/links"

	"github.com/gofiber/fiber/v3"
)

type StatsHandler struct {
	service *linkservice.Service
}

func NewStatsHandler(s *linkservice.Service) *StatsHandler {
	return &StatsHandler{service: s}
}

func (h *StatsHandler) Get(c fiber.Ctx) error {
	userID, ok := middleware.UserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":       "требуется авторизация",
			"redirect_to": "/auth",
		})
	}

	stats, err := h.service.GetStats(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "внутренняя ошибка сервера"})
	}

	c.Set("Cache-Control", "no-store")
	return c.JSON(stats)
}
