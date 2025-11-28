package main

import (
	"fmt"
	"os"
	"strings"

	"git-uts/config"
	"git-uts/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:      strings.Join(config.GetAllowedOrigins(), ","),
		AllowHeaders:      "Origin, Content-Type, Accept, Authorization",
		AllowMethods:      "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))


	config.MongoConnect()
	routes.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Not Found",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	fmt.Println("🚀 Server running on port", port)
	fmt.Println(app.Listen(":" + port))
}
