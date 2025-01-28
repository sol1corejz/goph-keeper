// Package main содержит реализацию основного сервера приложения Goph Keeper.
//
// Основные функции:
//   - Загрузка конфигурации сервера из файла с использованием паттерна Singleton.
//   - Установка подключения к базе данных.
//   - Инициализация HTTP-сервера с маршрутизацией для обработки запросов.
//   - Регистрация обработчиков для работы с пользователями, аутентификацией и управлением учетными данными.
//
// Используется фреймворк Fiber для обработки HTTP-запросов.
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
	internal "github.com/sol1corejz/goph-keeper/internal/server/handlers"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"sync"
)

// Singleton для конфига
var (
	config *configs.ServerConfig
	once   sync.Once
)

// main является точкой входа для запуска сервера.
// Он выполняет следующие шаги:
// 1. Загружает конфигурацию сервера из файла.
// 2. Подключается к базе данных.
// 3. Инициализирует приложение на основе Fiber и регистрирует маршруты.
// 4. Запускает сервер на указанном адресе.
func main() {
	var err error

	// Загрузка конфигурации
	config, err = LoadServerConfig("configs/server_config.yaml")
	if err != nil {
		log.Info("Failed to load server config", err.Error())
	}

	// Подключение к базе данных
	err = storage.ConnectDB(config)
	if err != nil {
		log.Info("Failed to connect to database", err.Error())
	}

	app := fiber.New()

	// Подключение Swagger UI с вашим YAML файлом
	app.Get("../../api/*", func(c *fiber.Ctx) error {
		c.SendFile("./docs/swagger.yaml")
		return nil
	})

	// Middleware для добавления конфигурации в контекст
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", config)
		return c.Next()
	})

	// Регистрация маршрутов
	app.Get("/", internal.RegisterHandler)
	app.Post("/register", internal.RegisterHandler)
	app.Post("/login", internal.LoginHandler)
	app.Post("/credentials", internal.AddCredentials)
	app.Get("/credentials", internal.GetCredentials)

	// Запуск сервера
	app.Listen(config.Server.Address)
}

// LoadServerConfig загружает конфигурацию сервера из указанного файла.
//
// Используется паттерн Singleton, чтобы гарантировать, что конфигурация загружается только один раз.
//
// Параметры:
//   - filePath: путь к YAML-файлу конфигурации.
//
// Возвращает:
//   - *configs.ServerConfig: объект конфигурации.
//   - error: ошибка, если не удалось загрузить конфигурацию.
func LoadServerConfig(filePath string) (*configs.ServerConfig, error) {
	var err error
	once.Do(func() {
		config, err = configs.LoadServerConfig(filePath)
	})
	return config, err
}
