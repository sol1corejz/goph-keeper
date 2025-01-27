package internal

// AuthPayload - входные данные для аутентификации
type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CredentialPayload - входные данные для записи данных
type CredentialPayload struct {
	Data string `json:"data"` // Основная информация (логин/пароль)
	Meta string `json:"meta"` // Метаинформация
}
