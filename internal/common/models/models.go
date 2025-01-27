// Package internal содержит внутреннюю логику приложэения
//   - Описание моделей
//   - Реализацию сервера
//   - Реализацию клиента
package internal

// Credential представляет данные учетной записи
type Credential struct {
	ID     string `json:"id"`      // Уникальный идентификатор
	UserID string `json:"user_id"` // Идентификатор пользователя
	Data   string `json:"data"`    // Основная информация (логин/пароль)
	Meta   string `json:"meta"`    // Метаинформация
}

// User представляет данные пользователя
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
