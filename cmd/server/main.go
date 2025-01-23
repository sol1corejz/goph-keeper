package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
	internal "github.com/sol1corejz/goph-keeper/internal/server/handlers"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"sync"
)

var (
	config *configs.ServerConfig
	once   sync.Once
)

func main() {
	var err error
	// Загрузка конфигурации
	config, err = LoadServerConfig("configs/server_config.yaml")
	if err != nil {
		log.Info("Failed to load server config", err.Error())
	}

	err = storage.ConnectDB(config)
	if err != nil {
		log.Info("Failed to connect to database", err.Error())
	}

	app := fiber.New()

	// Middleware для добавления конфигурации в контекст
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", config)
		return c.Next()
	})

	// Регистрация маршрутов
	app.Get("/", internal.RegisterHandler)
	app.Post("/register", internal.RegisterHandler)
	app.Post("/login", internal.LoginHandler)

	app.Listen(config.Server.Address)
}

// LoadServerConfig загружает конфиг с использованием Singleton
func LoadServerConfig(filePath string) (*configs.ServerConfig, error) {
	var err error
	once.Do(func() {
		config, err = configs.LoadServerConfig(filePath)
	})
	return config, err
}
