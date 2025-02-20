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
	userID string
)

var getCredentialsCmd = &cobra.Command{
	Use:   "get-credentials",
	Short: "Get credentials",
	Long:  "Получение данных пользователя",
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

		payloadData := &pb.GetCredentialsRequest{
			Token: token,
			Id:    userID,
		}

		resp, err := client.GetCredentials(ctx, payloadData)
		if err != nil {
			log.Fatalf("Ошибка получения данных: %v", err)
		}

		// Выводим ответ
		fmt.Println("Данные успешно получены!")
		fmt.Println(resp.Credentials)
		return

	},
}

func init() {
	rootCmd.AddCommand(getCredentialsCmd)
}
