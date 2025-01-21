package internal

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/client/auth"
	commonModels "github.com/sol1corejz/goph-keeper/internal/common/models"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func RegisterHandler(c *fiber.Ctx) error {
	//Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	var registerPayload internal.RegisterPayload

	//Парсинг входных данных
	err := json.Unmarshal(c.Body(), &registerPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//Генерация айди пользователя
	userUuid := uuid.New().String()

	//Хеширование пароля
	hashedPassword, err := HashPassword(registerPayload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	//Добавление пользователя в базу данных
	userData := commonModels.User{
		ID:       userUuid,
		Username: registerPayload.Username,
		Password: hashedPassword,
	}

	err = storage.DBStorage.CreateUser(userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//Генерация токена
	token, err := auth.GenerateToken(cfg, userUuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(auth.TokenExp),
		HTTPOnly: true,
	})
	//Отправка ответа
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": userData,
	})
}
