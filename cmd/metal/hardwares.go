package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"

	apiv1 "github.com/nicklasfrahm/infrastructure/api/v1"
)

const (
	repository      = "nicklasfrahm/infrastructure"
	fileURLTemplate = "https://raw.githubusercontent.com/%s/main/%s"
)

// Hardware returns a fiber app for the hardware resource.
func Hardware(useRemote bool) *fiber.App {
	app := fiber.New()

	app.Get("/:name<alpha>", func(c *fiber.Ctx) error {
		name := c.Params("name")

		var bytes []byte
		var err error

		path := fmt.Sprintf("configs/hardwares/%s.yaml", name)
		if useRemote {
			// Use the GitHub API to fetch the file.
			path = fmt.Sprintf(fileURLTemplate, repository, path)
			status, body, errs := fiber.Get(path).Bytes()
			if len(errs) > 0 {
				return errs[0]
			}
			if status != 200 {
				return fmt.Errorf("failed to fetch file: %s", path)
			}
			bytes = body
		} else {
			// Use the local file system to fetch the file.
			bytes, err = os.ReadFile(path)
			if err != nil {
				fmt.Printf("Failed to read file: %s", err)
				return err
			}
		}

		hardware := &apiv1.Hardware{}
		if err := yaml.Unmarshal(bytes, hardware); err != nil {
			return err
		}

		return c.Status(200).JSON(hardware)
	})

	return app
}
