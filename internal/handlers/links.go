package handlers

import (
	"errors"
	linksdto "goshortlinker/internal/data/dto/links"
	"goshortlinker/internal/middleware"
	linkservice "goshortlinker/internal/services/links"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
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
	data := fiber.Map{"IsAuthenticated": false}
	if user, ok := middleware.Current(c); ok {
		data["IsAuthenticated"] = true
		data["User"] = fiber.Map{"ID": user.UserID, "Name": user.Username}
	}
	return c.Render("home", data)
}

func (h *LinksHandler) GetRedirect(c fiber.Ctx) error {
	share_code := c.Params("share_code")

	link, err := h.service.GetLink(share_code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to retrieve link")
	}

	return c.Redirect().To(link.RedirectTo)
}

func (h *LinksHandler) Post(c fiber.Ctx) error {
	var payload linksdto.SaveRequest

	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid parse body")
	}

	userID := pgtype.Int4{}
	if id, ok := middleware.UserID(c); ok {
		userID = pgtype.Int4{Int32: id, Valid: true}
	}

	link, err := h.service.SaveAndShortenLink(userID, payload)
	if err != nil {
		if errors.Is(err, linkservice.ErrInvalidParams) {
			return c.Status(fiber.StatusBadRequest).SendString("некорректные параметры ссылки")
		}
		if errors.Is(err, linkservice.ErrSaveAndShortenLink) {
			return c.Status(fiber.StatusInternalServerError).SendString("could not shorten link")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}
	return c.Status(http.StatusCreated).JSON(link)
}
