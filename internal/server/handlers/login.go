package internal

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/goph-keeper/configs"
	auth "github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func LoginHandler(c *fiber.Ctx) error {
	//Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Переменная для входных данных
	var loginPayload internal.AuthPayload

	//Парсинг входных данных
	err := json.Unmarshal(c.Body(), &loginPayload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Получение пользователя из бд
	userData, err := storage.DBStorage.GetUser(loginPayload.Username)
	if err != nil {
		if errors.Is(storage.ErrNotFound, err) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong login or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Сравнение пароля из входных данных и из бд
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginPayload.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Wrong login or password",
		})
	}

	// Генерация токена
	token, err := auth.GenerateToken(cfg, userData.Password)

	// Установка токена в куки
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(auth.TokenExp),
		HTTPOnly: true,
	})

	// Отправка ответа
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": "login successfully",
	})
}
