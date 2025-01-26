package internal

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CredentialPayload struct {
	Data string `json:"data"` // Основная информация (логин/пароль)
	Meta string `json:"meta"` // Метаинформация
}
