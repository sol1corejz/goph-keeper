package internal

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// LoginHandler обрабатывает запросы на вход пользователя.
// Она парсит входные данные из тела запроса, проверяет логин и пароль пользователя
// с сохранёнными данными в базе данных, генерирует токен аутентификации и сохраняет его в cookie.
// В случае ошибки возвращается сообщение об ошибке с соответствующим статусом.
func LoginHandler(c *fiber.Ctx) error {
	// Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Переменная для входных данных
	var loginPayload internal.AuthPayload

	// Парсинг входных данных
	err := json.Unmarshal(c.Body(), &loginPayload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Получение пользователя из базы данных
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

	// Сравнение пароля из входных данных с паролем из базы данных
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginPayload.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Wrong login or password",
		})
	}

	// Генерация токена аутентификации
	token, err := auth.GenerateToken(cfg, userData.ID)

	// Установка токена в cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(auth.TokenExp),
		HTTPOnly: true,
	})

	// Отправка успешного ответа
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": "login successfully",
	})
}
