package app

import (
	"context"
	"goshortlinker/internal/handlers"
	genlinks "goshortlinker/internal/repos/gen/links"
	genusers "goshortlinker/internal/repos/gen/users"
	"goshortlinker/internal/services/links"
	users "goshortlinker/internal/services/registration"
	"goshortlinker/pkg/db/postgresql"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func ModuleDB() fx.Option {
	return fx.Provide(
		func() (*pgxpool.Pool, error) {
			return postgresql.InitDB(context.Background(), os.Getenv("POSTGRES_DSN"))
		},
	)
}

func ModuleRepositories() fx.Option {
	return fx.Provide(
		fx.Annotate(genlinks.New, fx.As(new(genlinks.Querier))),
		fx.Annotate(genusers.New, fx.As(new(genusers.Querier))),
	)
}

func ModuleServices() fx.Option {
	return fx.Provide(
		func(db *pgxpool.Pool, queries genlinks.Querier) *links.Service {
			return links.NewService(db, queries, 5*time.Second)
		},

		func(db *pgxpool.Pool, queries genusers.Querier) *users.Service {
			return users.NewService(db, queries, 5*time.Second)
		},
	)
}

func ModuleHandlers() fx.Option {
	return fx.Provide(
		handlers.NewLinksHandler,
		handlers.NewRegistrationHandler,
	)
}

func ModuleApp() fx.Option {
	return fx.Options(
		fx.Provide(NewApp),
		fx.Invoke(
			func(
				lc fx.Lifecycle,
				application *fiber.App,
				linksHandler *handlers.LinksHandler,
				registrationHandler *handlers.RegistrationHandler,
			) {
				RegisterRoutes(application, Handlers{
					LinksHandler:        linksHandler,
					RegistrationHandler: registrationHandler,
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
