package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
	internal "github.com/sol1corejz/goph-keeper/internal/common/models"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
)

// GetCredentials обрабатывает HTTP-запрос для получения списка пользовательских учетных данных.
//
// Запрос должен содержать:
//   - Cookie "token" с JWT токеном для авторизации.
//
// Алгоритм работы:
//  1. Извлечение конфигурации сервера из локального контекста.
//  2. Получение токена из cookie и проверка его наличия.
//  3. Проверка авторизации пользователя с помощью токена.
//  4. Получение учетных данных пользователя из хранилища.
//  5. Возврат списка учетных данных в формате JSON.
//
// Ответы:
//   - 200 OK: Учетные данные успешно извлечены и возвращены в ответе.
//   - 401 Unauthorized: Не предоставлен токен.
//   - 405 Method Not Allowed: Неверный токен.
//   - 500 Internal Server Error: Ошибка получения учетных данных.
//
// Параметры:
//   - c: Контекст запроса (Fiber).
//
// Возвращает:
//   - HTTP-ответ с соответствующим статусом и сообщением.
func GetCredentials(c *fiber.Ctx) error {
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

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(cfg, token)
	if err != nil {
		log.Info("token is invalid")
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "token is invalid",
		})
	}

	// Получение данных учетных записей из хранилища
	var credentialsData []internal.Credential
	credentialsData, err = storage.DBStorage.GetCredentials(userID)
	if err != nil {
		log.Info("failed to retrieve credentials")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve credentials",
		})
	}

	// Возврат списка учетных данных
	return c.JSON(fiber.Map{
		"credentials": credentialsData,
	})
}
