package app

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v3"
)

func NewApp() *fiber.App {
	fmt.Println("Initializing app...")
	engine := html.New("./views", ".html")

	return fiber.New(fiber.Config{
		Views: engine,
	})
}

func Run(app *fiber.App, port int) {
	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
