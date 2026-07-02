package app

import (
	"goshortlinker/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

type Hadnlers struct {
	HomeHandler *handlers.HomeHandler
}

func RegisterRoutes(app *fiber.App, handlers Hadnlers) {
	app.Get("/", handlers.HomeHandler.Get)
}
