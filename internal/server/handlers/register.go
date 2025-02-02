package internal

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// HashPassword принимает обычный пароль и возвращает хешированный пароль.
// Для хеширования используется bcrypt с дефолтным значением стоимости.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// RegisterHandler обрабатывает запросы на регистрацию пользователя.
// Она парсит входные данные из тела запроса, хеширует пароль пользователя,
// генерирует уникальный ID для пользователя и сохраняет данные пользователя в базе данных.
// Затем генерируется токен аутентификации, который сохраняется в cookie, и возвращается сообщение об успешной регистрации.
func RegisterHandler(c *fiber.Ctx) error {
	// Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Переменная для входных данных
	var registerPayload internal.AuthPayload

	// Парсинг входных данных
	err := json.Unmarshal(c.Body(), &registerPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Генерация айди пользователя
	userUuid := uuid.New().String()

	// Хеширование пароля
	hashedPassword, err := HashPassword(registerPayload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	// Добавление пользователя в базу данных
	userData := internal.User{
		ID:       userUuid,
		Username: registerPayload.Username,
		Password: hashedPassword,
	}

	// Создание пользователя в бд
	err = storage.DBStorage.CreateUser(userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Генерация токена
	token, err := auth.GenerateToken(cfg, userUuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	// Установка токена в куки
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(auth.TokenExp),
		HTTPOnly: true,
	})

	// Отправка ответа
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": "register successfully",
	})
}
