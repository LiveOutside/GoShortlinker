package app

import (
	"goshortlinker/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

type Handlers struct {
	LinksHandler *handlers.LinksHandler
	// AuthenticationHandler *handlers.AuthenticationHandler
}

func RegisterRoutes(app *fiber.App, handlers Handlers) {
	app.Get("/", handlers.LinksHandler.Get)
	app.Get("/:share_code", handlers.LinksHandler.GetRedirect)
	app.Post("/", handlers.LinksHandler.Post)

	// app.Get("/register", handlers.RegistrationHandler)
}
