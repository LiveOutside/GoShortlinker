package app

import (
	"context"
	"goshortlinker/internal/handlers"
	gencodes "goshortlinker/internal/repos/gen/activation_codes"
	genlinks "goshortlinker/internal/repos/gen/links"
	genusers "goshortlinker/internal/repos/gen/users"
	"goshortlinker/internal/services/links"
	activationcodes "goshortlinker/internal/services/mailer"
	users "goshortlinker/internal/services/registration"
	"goshortlinker/pkg/db/postgresql"
	mailer "goshortlinker/pkg/mailer"
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

func ModuleMailer() fx.Option {
	return fx.Provide(
		func() mailer.Mailer {
			return mailer.NewSMTPMailer(
				os.Getenv("SMTP_HOST"),
				os.Getenv("SMTP_PORT"),
				os.Getenv("SMTP_USER"),
				os.Getenv("SMTP_PASSWORD"),
				os.Getenv("SMTP_FROM"),
			)
		},
	)
}

func ModuleRepositories() fx.Option {
	return fx.Provide(
		fx.Annotate(genlinks.New, fx.As(new(genlinks.Querier))),
		fx.Annotate(genusers.New, fx.As(new(genusers.Querier))),
		fx.Annotate(gencodes.New, fx.As(new(gencodes.Querier))),
	)
}

func ModuleServices() fx.Option {
	return fx.Provide(
		func(db *pgxpool.Pool, queries genlinks.Querier) *links.Service {
			return links.NewService(db, queries, 5*time.Second)
		},

		func(db *pgxpool.Pool, queries gencodes.Querier, m mailer.Mailer) *activationcodes.Service {
			return activationcodes.NewService(db, queries, m, 5*time.Second)
		},

		func(db *pgxpool.Pool, queries genusers.Querier, activation *activationcodes.Service) *users.Service {
			return users.NewService(db, queries, activation, 5*time.Second)
		},
	)
}

func ModuleHandlers() fx.Option {
	return fx.Provide(
		handlers.NewLinksHandler,
		handlers.NewRegistrationHandler,
		handlers.NewMailerHandler,
		handlers.NewAuthenticationHandler,
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
				mailHandler *handlers.MailerHandler,
				authHandler *handlers.AuthenticationHandler,
			) {
				RegisterRoutes(application, Handlers{
					LinksHandler:          linksHandler,
					RegistrationHandler:   registrationHandler,
					MailHandler:           mailHandler,
					AuthenticationHandler: authHandler,
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
