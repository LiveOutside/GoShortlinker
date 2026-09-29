package app

import (
	"goshortlinker/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

type Handlers struct {
	LinksHandler        *handlers.LinksHandler
	RegistrationHandler *handlers.RegistrationHandler
	// AuthenticationHandler *handlers.AuthenticationHandler
}

func RegisterRoutes(app *fiber.App, handlers Handlers) {
	app.Get("/", handlers.LinksHandler.Get)
	app.Post("/", handlers.LinksHandler.Post)

	app.Get("/r/:share_code", handlers.LinksHandler.GetRedirect)

	app.Get("/register", handlers.RegistrationHandler.Get)
	app.Post("/register", handlers.RegistrationHandler.Post)
}
