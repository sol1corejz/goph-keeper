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
	pb "github.com/sol1corejz/goph-keeper/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	config *configs.ServerConfig
	once   sync.Once
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Инициализация конфигурации
	if err := initConfig(); err != nil {
		log.Fatal("Failed to load server config:", err)
	}

	// Подключение к базе данных
	if err := initDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Создание и настройка сервера
	app := setupServer()

	// Обработка сигнала завершения
	go handleShutdown(ctx, app, idleConnsClosed, sigint)

	// Запуск сервера
	startServer(app)

	// Запуск grpc
	grpcStart()
}

func initConfig() error {
	var err error
	config, err = LoadServerConfig("configs/server_config.yaml")
	return err
}

func initDatabase() error {
	return storage.ConnectDB(config)
}

func setupServer() *fiber.App {
	app := fiber.New()

	// Middleware для добавления конфигурации в контекст
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", config)
		return c.Next()
	})

	// Вывод информации о версии сборки
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)

	// Настройка маршрутов
	setupRoutes(app)

	return app
}

func setupRoutes(app *fiber.App) {
	app.Get("/", internal.RegisterHandler)
	app.Post("/register", internal.RegisterHandler)
	app.Post("/login", internal.LoginHandler)
	app.Post("/credentials", internal.AddCredentials)
	app.Post("/edit-credentials", internal.EditCredentials)
	app.Get("/credentials", internal.GetCredentials)
}

func handleShutdown(ctx context.Context, app *fiber.App, idleConnsClosed chan struct{}, sigint chan os.Signal) {
	<-sigint
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("HTTP server Shutdown failed:", err)
	}
	close(idleConnsClosed)
}

func startServer(app *fiber.App) {
	// Создание сертификата
	if !cert.CertExists() {
		log.Info("Generating new TLS certificate")
		certPEM, keyPEM := cert.GenerateCert()
		if err := cert.SaveCert(certPEM, keyPEM); err != nil {
			log.Errorf("failed to save TLS certificate: %v", err)
		}
	}

	log.Info("Loading existing TLS certificate")
	app.ListenTLS(config.Server.Address, cert.CertificateFilePath, cert.KeyFilePath)
}

func grpcStart() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Error(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterKeeperServer(s, &internal.KeeperServer{Config: config})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Error(err)
	}
}

func LoadServerConfig(filePath string) (*configs.ServerConfig, error) {
	var err error
	once.Do(func() {
		config, err = configs.LoadServerConfig(filePath)
	})
	return config, err
}
