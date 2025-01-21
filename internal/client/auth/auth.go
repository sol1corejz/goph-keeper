package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
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

var TokenExp = time.Hour * 60

func GenerateToken(config *configs.ServerConfig, userID string) (string, error) {

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
	return claims.UserID, nil
}

// CheckIsAuthorized проверяет наличие и валидность JWT токена в куках запроса.
// Возвращает UserID, если пользователь авторизован, или ошибку, если токен отсутствует или недействителен.
func CheckIsAuthorized(config *configs.ServerConfig, token string) (string, error) {
	// Извлекаем UserID из токена
	userID, err := ParseToken(config, token)
	if err != nil {
		return "", errors.New("token is invalid")
	}
	return userID, nil
}
