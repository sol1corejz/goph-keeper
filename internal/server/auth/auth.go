// Package auth содержит логику для работы с авторизацией.
//
// Функционал включает:
//   - Генерацию JWT токена.
//   - Парсинг JWT токена.
//   - Проверку валидности токена.
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
	// RegisteredClaims стандартные зарегистрированные поля JWT.
	jwt.RegisteredClaims
	// UserID уникальный идентификатор пользователя.
	UserID string
}

// TokenExp задает время жизни токена.
var TokenExp = time.Hour * 60

// GenerateToken создает JWT токен для указанного пользователя.
//
// Параметры:
//   - config: конфигурация сервера, содержащая секретный ключ для подписи токена.
//   - userID: уникальный идентификатор пользователя.
//
// Возвращает:
//   - строку, представляющую собой подписанный JWT токен.
//   - ошибку, если возникла проблема при создании токена.
func GenerateToken(config *configs.ServerConfig, userID string) (string, error) {

	log.Info(userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	secretKey := []byte(config.Security.JWTSecret)

	signedTokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedTokenString, nil
}

// ParseToken парсит и проверяет валидность указанного JWT токена.
//
// Параметры:
//   - config: конфигурация сервера, содержащая секретный ключ для проверки подписи токена.
//   - tokenString: строка токена, которую необходимо проверить.
//
// Возвращает:
//   - UserID, закодированный в токене, если он валиден.
//   - ошибку, если токен недействителен или содержит некорректные данные.
func ParseToken(config *configs.ServerConfig, tokenString string) (string, error) {
	claims := &Claims{}
	secretKey := []byte(config.Security.JWTSecret)

	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}

	// Проверяем валидность токена
	if !token.Valid {
		log.Info("Token is not valid")
		return "", errors.New("token is not valid")
	}

	// Возвращаем UserID из claims
	log.Info("Token is valid")

	// Проверка, что userID является валидным UUID
	if _, err = uuid.Parse(claims.UserID); err != nil {
		log.Info("UserID is not a valid UUID")
		return "", errors.New("userID in token is not valid")
	}

	return claims.UserID, nil
}

// CheckIsAuthorized проверяет наличие и валидность JWT токена.
//
// Параметры:
//   - config: конфигурация сервера, содержащая секретный ключ для проверки подписи токена.
//   - token: строка токена, которую необходимо проверить.
//
// Возвращает:
//   - UserID, если пользователь авторизован.
//   - ошибку, если токен отсутствует или недействителен.
func CheckIsAuthorized(config *configs.ServerConfig, token string) (string, error) {
	// Извлекаем UserID из токена
	userID, err := ParseToken(config, token)
	if err != nil {
		log.Info("Authorization failed:", err.Error())
		return "", errors.New("token is invalid")
	}

	return userID, nil
}
