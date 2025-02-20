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

// loginCmd представляет команду "login"
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Авторизация пользователя через gRPC",
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
		resp, err := client.Login(ctx, &pb.LoginRequest{
			UserData: userData,
		})
		if err != nil {
			if err != nil && strings.Contains(err.Error(), "not found") {
				fmt.Println("Неправильный логин или пароль!")
				return
			}
			log.Fatalf("Ошибка авторизации: %v", err)
		}

		// Сохраняем токен
		err = SaveTokenToFile(resp.Token)
		if err != nil {
			log.Fatalf("Ошибка сохранения токена: %v", err)
		}

		// Выводим ответ
		fmt.Println("Авторизация успешна!")
		return
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Добавляем флаги
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Имя пользователя")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль пользователя")

	// Флаги обязательны
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}
