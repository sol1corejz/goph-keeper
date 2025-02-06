package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/server/cert"
	internal "github.com/sol1corejz/goph-keeper/internal/server/handlers"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Singleton для конфигурации сервера
var (
	config *configs.ServerConfig
	once   sync.Once
)

// Глобальные переменные для информации о версии сборки.
var (
	buildVersion = "N/A" // Версия сборки, передается на этапе компиляции.
	buildDate    = "N/A" // Дата сборки, передается на этапе компиляции.
)

// main - входная точка приложения
func main() {
	// Канал сообщения о закртии соединения
	idleConnsClosed := make(chan struct{})
	// Канал для перенаправления прерываний
	sigint := make(chan os.Signal, 1)
	// Регистрация прерываний
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//Контекст отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Вывод информации о версии сборки.
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)

	// Горутина для обработки сигнала завершения
	go func() {
		<-sigint

		// Закрываем сервер
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Error("HTTP server Shutdown failed", err)
		}

		// Закрываем канал для уведомления о завершении
		close(idleConnsClosed)
	}()

	// Регистрация маршрутов
	app.Get("/", internal.RegisterHandler)
	app.Post("/register", internal.RegisterHandler)
	app.Post("/login", internal.LoginHandler)
	app.Post("/credentials", internal.AddCredentials)
	app.Post("/edit-credentials", internal.EditCredentials)
	app.Get("/credentials", internal.GetCredentials)

	// Создаем сертификат
	if !cert.CertExists() {
		log.Info("Generating new TLS certificate")
		certPEM, keyPEM := cert.GenerateCert()
		if err := cert.SaveCert(certPEM, keyPEM); err != nil {
			log.Errorf("failed to save TLS certificate: %w", err)
		}
	}

	log.Info("Loading existing TLS certificate")

	// Запускаем сервер
	app.ListenTLS(config.Server.Address, cert.CertificateFilePath, cert.KeyFilePath)

}

// LoadServerConfig загружает конфиг с использованием Singleton
func LoadServerConfig(filePath string) (*configs.ServerConfig, error) {
	var err error
	once.Do(func() {
		config, err = configs.LoadServerConfig(filePath)
	})
	return config, err
}
