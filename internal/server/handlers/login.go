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

// LoginHandler обрабатывает HTTP-запрос для авторизации пользователя.
//
// Запрос должен содержать:
//   - JSON-объект с полями "username" и "password".
//
// Алгоритм работы:
//  1. Извлечение конфигурации сервера из локального контекста.
//  2. Парсинг входных данных запроса.
//  3. Поиск пользователя в базе данных по имени.
//  4. Сравнение пароля из запроса с хешированным паролем из базы данных.
//  5. Генерация JWT токена при успешной авторизации.
//  6. Установка токена в HTTP cookie.
//  7. Возврат ответа о результате авторизации.
//
// Ответы:
//   - 202 Accepted: Успешная авторизация, токен установлен в cookie.
//   - 400 Bad Request: Некорректные входные данные.
//   - 401 Unauthorized: Неверный логин или пароль.
//   - 500 Internal Server Error: Ошибка при обработке запроса.
//
// Параметры:
//   - c: Контекст запроса (Fiber).
//
// Возвращает:
//   - HTTP-ответ с соответствующим статусом и сообщением.
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

	// Сравнение пароля из входных данных и хеша из базы данных
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginPayload.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Wrong login or password",
		})
	}

	// Генерация токена
	token, err := auth.GenerateToken(cfg, userData.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	// Установка токена в cookie
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
