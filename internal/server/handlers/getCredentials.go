package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
	internal "github.com/sol1corejz/goph-keeper/internal/common/models"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
)

func GetCredentials(c *fiber.Ctx) error {
	// Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Получение токане из куки
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

	var credentialsData []internal.Credential
	credentialsData, err = storage.DBStorage.GetCredentials(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve credentials",
		})
	}

	return c.JSON(fiber.Map{
		"credentials": credentialsData,
	})
}
