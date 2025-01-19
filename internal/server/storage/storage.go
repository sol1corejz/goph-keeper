package internal

import internal "github.com/sol1corejz/goph-keeper/internal/common/models"

type Storage interface {
	CreateUser(user internal.User) error
	GetUser(username, password string) (internal.User, error)
	SaveCredential(userID string, cred internal.Credential) error
	GetCredentials(userID string) ([]internal.Credential, error)
}
