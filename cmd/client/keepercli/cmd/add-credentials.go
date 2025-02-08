package cmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	pb "github.com/sol1corejz/goph-keeper/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// Флаги командной строки
var (
	data string
	meta string
)

var addCredentialsCmd = &cobra.Command{
	Use:   "add-credentials",
	Short: "Add credentials",
	Long:  "Добавление данных пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		// Устанавливаем соединение с gRPC сервером
		conn, err := grpc.NewClient("localhost:3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Ошибка подключения к gRPC: %v", err)
		}
		defer conn.Close()

		// Создаем клиента
		client := pb.NewKeeperClient(conn)

		// Формируем контекст с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// Получение токена авторизации
		token, err := ReadTokenFromFile()
		if err != nil {
			log.Fatalf("Ошибка получения токена: %v", err)
		}

		credentials := &pb.Credentials{
			Data: data,
			Meta: meta,
		}

		payloadData := &pb.AddCredentialsRequest{
			Token:       token,
			Credentials: credentials,
		}

		_, err = client.AddCredentials(ctx, payloadData)
		if err != nil {
			log.Fatalf("Ошибка добавления данных: %v", err)
		}

		// Выводим ответ
		fmt.Println("Данные успешно добавлены!")
		return

	},
}

func init() {
	rootCmd.AddCommand(addCredentialsCmd)

	// Добавляем флаги
	addCredentialsCmd.Flags().StringVarP(&data, "data", "d", "", "Данные пользователя")
	addCredentialsCmd.Flags().StringVarP(&meta, "meta", "m", "", "Метаданные")

	// Флаги обязательны
	addCredentialsCmd.MarkFlagRequired("data")
	addCredentialsCmd.MarkFlagRequired("meta")
}
