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

// AddCredentials обрабатывает HTTP-запрос для добавления новых пользовательских учетных данных.
//
// Запрос должен содержать:
//   - Cookie "token" с JWT токеном для авторизации.
//   - Тело запроса с данными учетных записей в формате JSON.
//
// Алгоритм работы:
//  1. Извлечение конфигурации сервера из локального контекста.
//  2. Получение токена из cookie и проверка его наличия.
//  3. Парсинг данных учетных записей из тела запроса.
//  4. Проверка авторизации пользователя с помощью токена.
//  5. Подготовка данных учетных записей для сохранения в базе данных.
//  6. Сохранение данных в хранилище.
//  7. Возврат статуса 201 Created при успешной операции.
//
// Формат тела запроса (пример):
//
//	{
//	  "data": "some data",
//	  "meta": "some meta information"
//	}
//
// Ответы:
//   - 201 Created: Учетные данные успешно добавлены.
//   - 401 Unauthorized: Не предоставлен токен.
//   - 405 Method Not Allowed: Неверный токен.
//   - 422 Unprocessable Entity: Ошибка парсинга данных запроса.
//   - 500 Internal Server Error: Ошибка сохранения учетных данных.
//
// Параметры:
//   - c: Контекст запроса (Fiber).
//
// Возвращает:
//   - HTTP-ответ с соответствующим статусом и сообщением.
func AddCredentials(c *fiber.Ctx) error {

	// Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Получение токена из cookie
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

	// Подготовка данных для БД
	credentialsData := commonModels.Credential{
		ID:     uuid.New().String(),
		UserID: userID,
		Data:   credentialsPayload.Data,
		Meta:   credentialsPayload.Meta,
	}

	// Сохранение данных в хранилище
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
