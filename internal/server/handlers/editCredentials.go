package internal

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
)

// EditCredentials обрабатывает запросы на редактирование учетных данных пользователя.
// Она извлекает токен из cookies, проверяет его валидность и авторизует пользователя.
// Затем она парсит входные данные, редактирует запись об учетных данных в базе данных и сохраняет её.
func EditCredentials(c *fiber.Ctx) error {

	// Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Получение токена из cookies
	token := c.Cookies("token")
	if token == "" {
		log.Info("No token cookie provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Парсинг входных данных
	var credentialsPayload internal.EditCredentialPayload
	err := json.Unmarshal(c.Body(), &credentialsPayload)
	if err != nil {
		log.Info("error unmarshalling credentials payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "failed to parse payload data",
		})
	}

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(cfg, token)
	if err != nil {
		log.Info("token is invalid")
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "token is invalid",
		})
	}

	// Подготовка данных для сохранения в базе данных
	credentialsData := internal.Credential{
		ID:     credentialsPayload.ID,
		UserID: userID,
		Data:   credentialsPayload.Data,
		Meta:   credentialsPayload.Meta,
	}

	// Сохранение учетных данных в базе данных
	err = storage.DBStorage.EditCredential(credentialsData)
	if err != nil {
		log.Info("failed to update credential")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update credential data",
		})
	}

	// Отправка успешного ответа
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": "credentials updated",
	})
}
