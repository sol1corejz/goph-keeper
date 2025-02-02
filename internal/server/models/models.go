// Package internal содержит все бизнес-логики приложения и обеспечивает взаимодействие с
// хранилищем данных, а также авторизацией и обработкой запросов.
// В этом пакете находятся модели данных, обработчики HTTP-запросов, а также логика,
// отвечающая за взаимодействие с базой данных и обработку аутентификации пользователей.
package internal

// AuthPayload представляет данные для аутентификации пользователя.
// Используется при отправке запроса на вход в систему.
type AuthPayload struct {
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль
}

// CredentialPayload содержит учетные данные пользователя.
// Используется для передачи логина и пароля с дополнительными метаданными.
type CredentialPayload struct {
	Data string `json:"data"` // Основная информация (например, логин и пароль)
	Meta string `json:"meta"` // Дополнительные метаданные (например, описание, время создания)
}

// Credential представляет учетную запись пользователя, сохраненную в системе.
type Credential struct {
	ID     string `json:"id"`      // Уникальный идентификатор учетной записи
	UserID string `json:"user_id"` // Идентификатор владельца учетной записи
	Data   string `json:"data"`    // Основная информация (логин/пароль)
	Meta   string `json:"meta"`    // Дополнительные метаданные
}

// User представляет зарегистрированного пользователя в системе.
type User struct {
	ID       string `json:"id"`       // Уникальный идентификатор пользователя
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Хешированный пароль пользователя
}
