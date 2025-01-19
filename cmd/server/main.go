package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
)

func main() {
	config, err := configs.LoadServerConfig("configs/server_config.yaml")
	if err != nil {
		log.Info("Failed to load server config", err.Error())
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		return c.SendString("register")
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("login")
	})

	app.Listen(config.Server.Address)
}
