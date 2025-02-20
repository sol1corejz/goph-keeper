package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
)

// GetCredentials обрабатывает запросы на получение учетных данных пользователя.
// Она извлекает токен из cookies, проверяет его валидность и авторизует пользователя.
// После этого она извлекает учетные данные из базы данных и возвращает их в ответе.
func GetCredentials(c *fiber.Ctx) error {
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

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(cfg, token)
	if err != nil {
		log.Info("token is invalid")
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "token is invalid",
		})
	}

	// Получение учетных данных пользователя из базы данных
	var credentialsData []internal.Credential
	credentialsData, err = storage.DBStorage.GetCredentials(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve credentials",
		})
	}

	// Отправка учетных данных в ответе
	return c.JSON(fiber.Map{
		"credentials": credentialsData,
	})
}
