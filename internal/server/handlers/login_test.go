package internal_test

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func LoginHandler(c *fiber.Ctx) error {
	// Здесь логика аутентификации (например, проверка данных)
	var payload internal.AuthPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Проверка данных пользователя
	if payload.Username == "testuser" && payload.Password == "correctpassword" {
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Login successful"})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
}

func TestLoginHandler(t *testing.T) {
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
			name: "Test successful login",
			inputBody: internal.AuthPayload{
				Username: "testuser",
				Password: "correctpassword",
			},
			want: want{
				code: fiber.StatusAccepted,
				body: "Login successful",
			},
		},
		{
			name: "Test failed login",
			inputBody: internal.AuthPayload{
				Username: "testuser",
				Password: "wrongpassword",
			},
			want: want{
				code: fiber.StatusUnauthorized,
				body: "Invalid credentials",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем новое приложение Fiber
			app := fiber.New()

			// Регистрируем хендлер
			app.Post("/login", LoginHandler)

			// Создаем тело запроса
			body, _ := json.Marshal(test.inputBody)

			// Создаем запрос
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
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
