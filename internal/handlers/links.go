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

func (h *LinksHandler) renderLinkError(c fiber.Ctx, status int, title, message string) error {
	return c.Status(status).Render("link_error", fiber.Map{
		"Title":   title,
		"Message": message,
	})
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

	link, err := h.service.GetLinkForRedirect(share_code)
	if err != nil {
		switch {
		case errors.Is(err, linkservice.ErrLinkNotFound):
			return h.renderLinkError(c, fiber.StatusNotFound,
				"Link not found",
				"The link you are trying to access does not exist or has been removed.")
		// REWORK?: maybe need to return different error types instead of ErrLinkUnavailable
		case errors.Is(err, linkservice.ErrLinkUnavailable):
			return h.renderLinkError(c, fiber.StatusBadRequest,
				"Link unavailable",
				"The link you are trying to access is currently unavailable.")
		default:
			return h.renderLinkError(c, fiber.StatusInternalServerError,
				"Something went wrong",
				"An unexpected error occurred while processing your request.")
		}
	}

	return c.Render("redirect", fiber.Map{
		"ShareCode":     share_code,
		"RedirectTimer": link.RedirectTimer,
	})
}

func (h *LinksHandler) ConfirmRedirect(c fiber.Ctx) error {
	response, err := h.service.ConfirmRedirect(c.Params("share_code"), c.IP())
	if err != nil {
		switch {
		case errors.Is(err, linkservice.ErrLinkNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "link not found"})
		case errors.Is(err, linkservice.ErrLinkUnavailable):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "link is unavailable"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
	}
	return c.JSON(response)
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
