package main

// Create a basic gofiber v2 app that exposes a single health endpoint.
// The health endpoint will return a 200 status code and a simple message.
// The health endpoint will be exposed on the environment variable PORT.

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

var version = "dev"

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Important: Avoid using "/healthz" as this is reserved by Google Cloud Run.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  http.StatusOK,
			"title":   http.StatusText(http.StatusOK),
			"version": version,
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status": http.StatusNotFound,
			"title":  http.StatusText(http.StatusNotFound),
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
