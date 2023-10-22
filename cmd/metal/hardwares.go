package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"

	apiv1 "github.com/nicklasfrahm/infrastructure/api/v1"
)

const (
	repoOwner          = "nicklasfrahm"
	repoName           = "infrastructure"
	hardwareConfigPath = "configs/hardwares"
	fileURLTemplate    = "https://raw.githubusercontent.com/%s/%s/main/%s"
)

// Hardware returns a fiber app for the hardware resource.
func Hardware(useRemote bool) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		hardwares, err := ListHardwares()
		if err != nil {
			return err
		}

		return c.Status(200).JSON(hardwares)
	})

	app.Get("/:name<alpha>", func(c *fiber.Ctx) error {
		name := c.Params("name")

		var bytes []byte
		var err error

		path := fmt.Sprintf("%s/%s.yaml", hardwareConfigPath, name)
		if useRemote {
			// Use the GitHub API to fetch the file.
			path = fmt.Sprintf(fileURLTemplate, repoOwner, repoName, path)
			status, body, errs := fiber.Get(path).Bytes()
			if len(errs) > 0 {
				return errs[0]
			}
			if status != 200 {
				if status == 404 {
					return c.Next()
				}
				return fmt.Errorf("failed to fetch file: %s", path)
			}
			bytes = body
		} else {
			// Use the local file system to fetch the file.
			bytes, err = os.ReadFile(path)
			if err != nil {
				return err
			}
		}

		// Parse the YAML file into a hardware object.
		hardware := &apiv1.Hardware{}
		if err := yaml.Unmarshal(bytes, hardware); err != nil {
			return err
		}

		return c.Status(200).JSON(hardware)
	})

	return app
}

// ListHardwares returns a list of all available hardwares.
func ListHardwares() ([]*apiv1.Hardware, error) {
	client := NewGitHubClient()

	_, dir, _, err := client.Repositories.GetContents(context.TODO(), repoOwner, repoName, hardwareConfigPath, nil)
	if err != nil {
		return nil, err
	}

	hardwares := []*apiv1.Hardware{}
	for _, entry := range dir {
		file, _, _, err := client.Repositories.GetContents(context.TODO(), repoOwner, repoName, entry.GetPath(), nil)
		if err != nil {
			return nil, err
		}

		content, err := file.GetContent()
		if err != nil {
			return nil, err
		}

		hardware := &apiv1.Hardware{}
		if err := yaml.Unmarshal([]byte(content), hardware); err != nil {
			return nil, err
		}

		hardwares = append(hardwares, hardware)
	}

	return hardwares, nil
}

// ReadHardware returns a hardware by name.
func ReadHardware(name string) (*apiv1.Hardware, error) {
	return nil, nil
}
