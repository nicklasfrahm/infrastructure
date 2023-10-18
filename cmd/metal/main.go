package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

var version = "dev"

func init() {
	// Allow usage of JSON logs in production.
	format := os.Getenv("LOG_FORMAT")
	if format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	level := os.Getenv("LOG_LEVEL")
	log.SetLevel(log.InfoLevel)
	if level != "" {
		logLevel, err := log.ParseLevel(level)
		if err != nil {
			log.Fatalf("Failed to parse log level: %s", err)
		}
		log.SetLevel(logLevel)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          ErrorHandler,
	})

	var useRemote bool
	if _, err := os.Stat("configs"); os.IsNotExist(err) {
		useRemote = true
	}

	// Configure resource routes.
	app.Mount("/hardwares", Hardware(useRemote))

	// IMPORTANT: Avoid using "/healthz" as this is reserved by Google Cloud Run.
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

	log.Infof("Starting server: http://0.0.0.0:%s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

// ErrorHandler is a custom error handler for fiber.
func ErrorHandler(c *fiber.Ctx, err error) error {
	if os.IsNotExist(err) {
		return c.Next()
	}

	log.Error(err)

	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"status": code,
		"title":  http.StatusText(code),
	})
}
