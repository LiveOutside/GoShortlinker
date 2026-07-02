package app

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/template/html/v3"
)

func NewApp() *fiber.App {
	fmt.Println("Initializing app...")
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET,POST,PUT,DELETE,OPTIONS"},
		AllowHeaders: []string{"Origin, Content-Type, Accept"},
	}))

	return app
}

func Run(app *fiber.App, port int) {
	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
