package internal_test

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func RegisterHandler(c *fiber.Ctx) error {
	// Получение конфига из контекста
	cfg := c.Locals("config").(*configs.ServerConfig)

	// Переменная для входных данных
	var registerPayload internal.AuthPayload

	// Парсинг входных данных
	err := json.Unmarshal(c.Body(), &registerPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if registerPayload.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username is required",
		})
	}

	if registerPayload.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password is required",
		})
	}

	// Генерация айди пользователя
	userUuid := uuid.New().String()

	// Генерация токена
	token, err := auth.GenerateToken(cfg, userUuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
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
		"message": "register successfully",
	})
}

func TestRegisterHandler(t *testing.T) {
	type want struct {
		code int
		body string
	}

	tests := []struct {
		name      string
		inputBody internal.AuthPayload
		want      want
	}{
		{
			name: "Test successful registration",
			inputBody: internal.AuthPayload{
				Username: "newuser",
				Password: "newpassword123",
			},
			want: want{
				code: fiber.StatusCreated,
				body: "register successfully",
			},
		},
		{
			name: "Test failed registration due to invalid input",
			inputBody: internal.AuthPayload{
				Username: "newuser",
				Password: "",
			},
			want: want{
				code: fiber.StatusBadRequest,
				body: "Password is required",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем новое приложение Fiber
			app := fiber.New()

			config, err := configs.LoadServerConfig("../../../configs/server_config.yaml")
			if err != nil {
				t.Errorf("could not load server config")
			}

			// Middleware для добавления конфигурации в контекст
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("config", config)
				return c.Next()
			})

			// Регистрируем хендлер
			app.Post("/register", RegisterHandler)

			// Создаем тело запроса
			body, _ := json.Marshal(test.inputBody)

			// Создаем запрос
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Отправляем запрос через Fiber
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Проверка ответа
			assert.Equal(t, test.want.code, resp.StatusCode)

			// Чтение и проверка тела ответа
			var resBody map[string]string
			json.NewDecoder(resp.Body).Decode(&resBody)
			assert.Equal(t, test.want.body, resBody["message"])
		})
	}
}
