package cmd

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
	"time"

	pb "github.com/sol1corejz/goph-keeper/proto"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

// Флаги командной строки
var (
	username string
	password string
)

// registerCmd представляет команду "register"
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация нового пользователя через gRPC",
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

		userData := &pb.User{
			Username: username,
			Password: password,
		}
		// Отправляем запрос
		resp, err := client.Register(ctx, &pb.RegisterRequest{
			UserData: userData,
		})
		if err != nil {
			if err != nil && strings.Contains(err.Error(), "already exists") {
				fmt.Println("Пользователь уже зарегестрирован!")
				return
			}
			log.Fatalf("Ошибка регистрации: %v", err)
			return
		}

		// Сохраняем токен
		err = SaveTokenToFile(resp.Token)
		if err != nil {
			log.Fatalf("Ошибка сохранения токена: %v", err)
		}

		// Выводим ответ
		fmt.Println("Регистрация успешна!")
		return
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Добавляем флаги
	registerCmd.Flags().StringVarP(&username, "username", "u", "", "Имя пользователя")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль пользователя")

	// Флаги обязательны
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
}
