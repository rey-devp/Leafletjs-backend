package routes

import (
	"git-uts/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Get("/locations", handler.GetLocations)
	api.Post("/locations", handler.CreateLocation)
	api.Put("/locations/:id", handler.UpdateLocation)
	api.Delete("/locations/:id", handler.DeleteLocation)
}