package main

import (
	"goshortlinker/internal/app"

	"go.uber.org/fx"
)

// @title GoShortlinker
// @version 1.0
// @description API for shortening links
// @host localhost:8080
// BasePath /
func main() {
	fx.New(
		app.ModuleRepositories(),
		app.ModuleServices(),
		app.ModuleHandlers(),
		app.ModuleApp(),
	).Run()
}
