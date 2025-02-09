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

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if err := initConfig(); err != nil {
		log.Fatal("Failed to load server config:", err)
	}

	if err := initDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	app := setupServer()

	// Запускаем HTTP сервер в отдельной горутине
	go startServer(app)

	// Запускаем gRPC сервер в отдельной горутине
	grpcClosed := make(chan struct{})
	go grpcStart(ctx, grpcClosed)

	// Ожидание сигнала завершения
	<-sigint
	log.Info("Получен сигнал завершения, останавливаем серверы...")

	cancel() // Отправляем сигнал завершения контексту

	// Завершаем Fiber
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Ошибка при завершении HTTP сервера:", err)
	}

	// Ждём завершения gRPC-сервера
	<-grpcClosed
	log.Info("Сервер полностью завершён")
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
	if err := app.ListenTLS(config.Server.Address, cert.CertificateFilePath, cert.KeyFilePath); err != nil {
		log.Fatal("Ошибка запуска HTTP сервера:", err)
	}
}

func grpcStart(ctx context.Context, closed chan struct{}) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Error("Ошибка при запуске gRPC сервера:", err)
		close(closed)
		return
	}

	s := grpc.NewServer()
	pb.RegisterKeeperServer(s, &internal.KeeperServer{Config: config})

	go func() {
		<-ctx.Done()
		log.Info("Останавливаем gRPC сервер...")
		s.GracefulStop()
		close(closed)
	}()

	log.Info("gRPC сервер запущен на порту 3200")
	if err := s.Serve(listen); err != nil {
		log.Error("Ошибка работы gRPC сервера:", err)
	}
}

// LoadServerConfig - Функция загрузки конфигурации из файла с конфигурациями
func LoadServerConfig(filePath string) (*configs.ServerConfig, error) {
	var err error
	once.Do(func() {
		config, err = configs.LoadServerConfig(filePath)
	})
	return config, err
}
