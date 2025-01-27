package internal

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sol1corejz/goph-keeper/configs"
	commonModels "github.com/sol1corejz/goph-keeper/internal/common/models"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// HashPassword хеширует пароль с использованием bcrypt.
// Принимает:
//   - password: строка, представляющая пароль для хеширования.
//
// Возвращает:
//   - строка: хешированный пароль.
//   - error: ошибка, если процесс хеширования завершился неудачно.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// RegisterHandler обрабатывает запросы на регистрацию новых пользователей.
//
// Запрос должен содержать:
//   - JSON-объект с полями "username" и "password".
//
// Алгоритм работы:
//  1. Извлечение конфигурации сервера из локального контекста.
//  2. Парсинг входных данных запроса.
//  3. Генерация уникального идентификатора пользователя.
//  4. Хеширование пароля.
//  5. Добавление нового пользователя в базу данных.
//  6. Генерация JWT токена для авторизованного сеанса.
//  7. Установка токена в HTTP cookie.
//  8. Возврат ответа о результате регистрации.
//
// Ответы:
//   - 201 Created: Успешная регистрация пользователя, токен установлен в cookie.
//   - 400 Bad Request: Некорректные входные данные.
//   - 500 Internal Server Error: Ошибка при обработке запроса.
//
// Параметры:
//   - c: Контекст запроса (Fiber).
//
// Возвращает:
//   - HTTP-ответ с соответствующим статусом и сообщением.
func RegisterHandler(c *fiber.Ctx) error {
	// Получение конфигурации из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Переменная для входных данных
	var registerPayload internal.AuthPayload

	// Парсинг входных данных
	err := json.Unmarshal(c.Body(), &registerPayload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Генерация уникального идентификатора пользователя
	userUuid := uuid.New().String()

	// Хеширование пароля
	hashedPassword, err := HashPassword(registerPayload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to hash password",
		})
	}

	// Создание записи пользователя для базы данных
	userData := commonModels.User{
		ID:       userUuid,
		Username: registerPayload.Username,
		Password: hashedPassword,
	}

	// Добавление пользователя в базу данных
	err = storage.DBStorage.CreateUser(userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Генерация JWT токена
	token, err := auth.GenerateToken(cfg, userUuid)
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

	// Отправка успешного ответа
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": "register successfully",
	})
}
