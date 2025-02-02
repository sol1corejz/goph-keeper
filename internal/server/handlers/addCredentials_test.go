package internal_test

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkIsAuthorized(token string) (string, error) {
	if token == "valid-token" {
		return "123", nil
	} else if token == "invalid-token" {
		return "", errors.New("invalid token")
	} else {
		return "", errors.New("token is invalid")
	}
}

func AddCredentials(c *fiber.Ctx) error {
	// Получение токена из cookies
	token := c.Cookies("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	// Проверка авторизации
	_, err := checkIsAuthorized(token)
	if err != nil {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"message": "token is invalid",
		})
	}

	// Отправка успешного ответа
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "credential added",
	})
}

func TestAddCredentialsHandler(t *testing.T) {
	type want struct {
		code int
		body string
	}

	tests := []struct {
		name  string
		token string
		body  string
		want  want
	}{
		{
			name:  "Test successful credential addition",
			token: "valid-token",
			want: want{
				code: fiber.StatusCreated,
				body: "credential added",
			},
		},
		{
			name:  "Test missing token",
			token: "",
			want: want{
				code: fiber.StatusUnauthorized,
				body: "unauthorized",
			},
		},
		{
			name:  "Test invalid token",
			token: "invalid-token",
			want: want{
				code: fiber.StatusMethodNotAllowed,
				body: "token is invalid",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем новое приложение Fiber
			app := fiber.New()

			// Регистрируем хендлер
			app.Post("/credentials", AddCredentials)

			// Создаем запрос с токеном в cookies
			req := httptest.NewRequest(http.MethodPost, "/credentials", nil)
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: test.token,
			})

			// Отправляем запрос через Fiber
			resp, _ := app.Test(req)

			// Проверка ответа
			assert.Equal(t, test.want.code, resp.StatusCode)

			// Чтение и проверка тела ответа
			var resBody map[string]string
			json.NewDecoder(resp.Body).Decode(&resBody)
			assert.Equal(t, test.want.body, resBody["message"])
		})
	}
}
