package handlers

import (
	"errors"
	linksdto "goshortlinker/internal/data/dto/links"
	linkservice "goshortlinker/internal/services/links"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type LinksHandler struct {
	service *linkservice.Service
}

func NewLinksHandler(s *linkservice.Service) *LinksHandler {
	return &LinksHandler{
		service: s,
	}
}

// @Summary Get Short Link
// @Description Get short link with specified parameters
// @Tags home
// @Produce html
// @Success 200 {string} string "HTML page with shortlink service"
// @Failure 500 {string} string "Internal Server Error"
// @Router / [get]
func (h *LinksHandler) Get(c fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}

func (h *LinksHandler) Post(c fiber.Ctx) error {
	var payload linksdto.SaveRequest

	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid parse body")
	}

	link, err := h.service.SaveAndShortenLink(payload)
	if err != nil {
		if errors.Is(err, linkservice.ErrSaveAndShortenLink) {
			return c.Status(fiber.StatusInternalServerError).SendString("could not shorten link")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}
	return c.Status(http.StatusCreated).JSON(link)
}
