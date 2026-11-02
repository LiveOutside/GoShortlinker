package app

import (
	"goshortlinker/internal/handlers"
	"goshortlinker/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type Handlers struct {
	LinksHandler          *handlers.LinksHandler
	RegistrationHandler   *handlers.RegistrationHandler
	MailHandler           *handlers.MailerHandler
	AuthenticationHandler *handlers.AuthenticationHandler
	ProfileHandler        *handlers.ProfileHandler
	StatsHandler          *handlers.StatsHandler
	SocialHandler         *handlers.SocialHandler
}

func RegisterRoutes(app *fiber.App, handlers Handlers) {
	loginLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "слишком много попыток, попробуйте позже",
			})
		},
	})

	guestOnly := middleware.RedirectIfAuthenticated()

	app.Get("/", handlers.LinksHandler.Get)
	app.Post("/", handlers.LinksHandler.Post)

	app.Get("/r/:share_code", handlers.LinksHandler.GetRedirect)
	app.Post("/r/:share_code/confirm", handlers.LinksHandler.ConfirmRedirect)

	app.Get("/register", guestOnly, handlers.RegistrationHandler.Get)
	app.Post("/register", handlers.RegistrationHandler.Post)

	app.Get("/activate", guestOnly, handlers.MailHandler.Get)
	app.Post("/activate", handlers.MailHandler.Post)
	app.Post("/activate/resend", handlers.MailHandler.Resend)

	app.Get("/auth", guestOnly, handlers.AuthenticationHandler.Get)
	app.Post("/auth", loginLimiter, handlers.AuthenticationHandler.Post)
	app.Post("/logout", handlers.AuthenticationHandler.Logout)

	app.Get("/auth/:provider/callback", handlers.SocialHandler.Callback)
	app.Get("/auth/:provider", guestOnly, handlers.SocialHandler.StartLogin)

	app.Get("/profile", handlers.ProfileHandler.Get)
	app.Put("/profile", handlers.ProfileHandler.Update)
	app.Get("/profile/stats", handlers.StatsHandler.Get)
	app.Get("/profile/social/:provider/link", handlers.SocialHandler.StartLink)
	app.Post("/profile/social/:provider/unlink", handlers.SocialHandler.Unlink)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{
			"Title":   "Page not found",
			"Message": "The page you are trying to access does not exist.",
		})
	})
}
