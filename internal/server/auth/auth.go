// Package auth содержит логику для генерации и проверки JSON Web Token (JWT),
// а также для авторизации пользователей в приложении.
// В этом пакете определены структуры для токенов, функции для создания,
// парсинга и валидации токенов, а также проверка авторизации пользователей
// на основе переданных токенов.
package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sol1corejz/goph-keeper/configs"
	"time"
)

// Claims структура, содержащая информацию о пользователе,
// которая будет закодирована в JWT токене.
type Claims struct {
	// Зарегистрированные стандартные поля JWT.
	jwt.RegisteredClaims
	// UserID - уникальный идентификатор пользователя.
	UserID string
}

// TokenExp - время жизни JWT токена.
var TokenExp = time.Hour * 60

// GenerateToken генерирует новый JWT токен для пользователя.
// Включает в токен информацию о времени истечения и UserID.
func GenerateToken(config *configs.ServerConfig, userID string) (string, error) {

	log.Info(userID)

	// Создание нового токена с указанными claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	// Получение секретного ключа из конфигурации
	secretKey := []byte(config.Security.JWTSecret)

	// Подписание токена с использованием секрета
	signedTokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Возврат подписанного токена
	return signedTokenString, nil
}

// ParseToken парсит JWT токен и извлекает из него UserID.
// Возвращает UserID, если токен действителен, или ошибку в случае неудачи.
func ParseToken(config *configs.ServerConfig, tokenString string) (string, error) {
	claims := &Claims{}
	secretKey := []byte(config.Security.JWTSecret)

	// Парсинг токена с извлечением данных в claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}

	// Проверка валидности токена
	if !token.Valid {
		log.Info("Token is not valid")
		return "", errors.New("token is not valid")
	}

	// Возврат UserID из claims
	log.Info("Token is valid")

	// Проверка, что userID является валидным UUID
	if _, err = uuid.Parse(claims.UserID); err != nil {
		log.Info("UserID is not a valid UUID")
		return "", errors.New("userID in token is not valid")
	}

	return claims.UserID, nil
}

// CheckIsAuthorized проверяет наличие и валидность JWT токена в куках запроса.
// Возвращает UserID, если пользователь авторизован, или ошибку, если токен отсутствует или недействителен.
func CheckIsAuthorized(config *configs.ServerConfig, token string) (string, error) {
	// Извлекаем UserID из токена
	userID, err := ParseToken(config, token)
	if err != nil {
		log.Info("Authorization failed:", err.Error())
		return "", errors.New("token is invalid")
	}

	return userID, nil
}
