package app

import (
	"context"
	"goshortlinker/internal/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"go.uber.org/fx"
)

func ModuleHandlers() fx.Option {
	return fx.Provide(
		handlers.NewHomeHandler,
	)
}

func ModuleApp() fx.Option {
	return fx.Options(
		fx.Provide(NewApp),
		fx.Invoke(
			func(
				lc fx.Lifecycle,
				application *fiber.App,
				homeHandler *handlers.HomeHandler,
			) {
				RegisterRoutes(application, Hadnlers{
					HomeHandler: homeHandler,
				})

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := application.Listen(":8080"); err != nil {
								log.Errorf("Falied to start listening (fiber): %w", err)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return application.ShutdownWithContext(ctx)
					},
				})
			},
		),
	)
}
