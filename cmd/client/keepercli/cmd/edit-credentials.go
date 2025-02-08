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

var dataID string

var editCredentialsCmd = &cobra.Command{
	Use:   "edit-credentials",
	Short: "Edit credentials",
	Long:  "Обновление данных пользователя",
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

		payloadData := &pb.EditCredentialsRequest{
			Id:          dataID,
			Credentials: credentials,
			Token:       token,
		}

		_, err = client.EditCredentials(ctx, payloadData)
		if err != nil {
			log.Fatalf("Ошибка обновления данных: %v", err)
		}

		// Выводим ответ
		fmt.Println("Данные успешно обновлены!")
		return

	},
}

func init() {
	rootCmd.AddCommand(editCredentialsCmd)

	// Добавляем флаги
	editCredentialsCmd.Flags().StringVarP(&dataID, "id", "i", "", "Идентификатор данных")
	editCredentialsCmd.Flags().StringVarP(&data, "data", "d", "", "Данные пользователя")
	editCredentialsCmd.Flags().StringVarP(&meta, "meta", "m", "", "Метаданные")

	// Флаги обязательны
	editCredentialsCmd.MarkFlagRequired("id")
	editCredentialsCmd.MarkFlagRequired("data")
	editCredentialsCmd.MarkFlagRequired("meta")
}
