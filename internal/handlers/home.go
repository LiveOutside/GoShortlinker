package handlers

import "github.com/gofiber/fiber/v3"

type HomeHandler struct {
	service *homeService.Service
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// @Summary Get Short Link
// @Description Get short link with specified parameters
// @Tags home
// @Produce html
// @Success 200 {string} string "HTML page with shortlink service"
// @Failure 500 {string} string "Internal Server Error"
// @Router / [get]
func (h *HomeHandler) Get(c fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}
