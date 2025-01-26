package internal

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/sol1corejz/goph-keeper/configs"
	commonModels "github.com/sol1corejz/goph-keeper/internal/common/models"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
)

func AddCredentials(c *fiber.Ctx) error {

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

	// Парсинг входных данных
	var credentialsPayload internal.CredentialPayload
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

	// Подготовка данных для бд
	credentialsData := commonModels.Credential{
		ID:     uuid.New().String(),
		UserID: userID,
		Data:   credentialsPayload.Data,
		Meta:   credentialsPayload.Meta,
	}

	err = storage.DBStorage.SaveCredential(credentialsData)
	if err != nil {
		log.Info("failed to save credential")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save credential data",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": "credential added",
	})

}
