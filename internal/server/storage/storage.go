package internal

import (
	"database/sql"
	"errors"
	log "github.com/gofiber/fiber/v2/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sol1corejz/goph-keeper/configs"
	internal "github.com/sol1corejz/goph-keeper/internal/common/models"
)

// Storage - Абстракия для хранилища, для возможного расширения и использования разных видов хранения
// Также нужно для написания тестов, чтобы содавать моковое хранилище
type Storage interface {
	CreateUser(user internal.User) error
	GetUser(username string) (internal.User, error)
	SaveCredential(userID string, cred internal.Credential) error
	GetCredentials(userID string) ([]internal.Credential, error)
}

// StorageImpl структура реализующая интерфейс Storage
type StorageImpl struct {
	DB *sql.DB
}

// DBStorage объект использующийся для использования методов хранилища
var DBStorage StorageImpl
var ErrNotFound = errors.New("not found")

func ConnectDB(cfg *configs.ServerConfig) error {

	if cfg.Storage.ConnectionString == "" {
		return errors.New("no connection string provided")
	}

	db, err := sql.Open("pgx", cfg.Storage.ConnectionString)
	if err != nil {
		log.Fatal(err)
		return err
	}

	DBStorage.DB = db

	// Создание таблицы users
	_, err = DBStorage.DB.Query(`
		CREATE TABLE IF NOT EXISTS users (
    		uuid UUID PRIMARY KEY,
    		username TEXT NOT NULL UNIQUE,
    		password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (s *StorageImpl) CreateUser(user internal.User) error {

	_, err := DBStorage.DB.Exec(`
		INSERT INTO users (uuid, username, password) VALUES ($1, $2, $3)
	`, user.ID, user.Username, user.Password)

	if err != nil {
		log.Info("failed to create user")
		return err
	}

	return nil
}
func (s *StorageImpl) GetUser(username string) (internal.User, error) {

	var user internal.User
	err := DBStorage.DB.QueryRow(`
		SELECT * FROM users WHERE username=$1
	`, username).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.User{}, ErrNotFound
		}
		return internal.User{}, err
	}

	return user, nil
}
func (s *StorageImpl) SaveCredential(userID string, cred internal.Credential) error {
	return nil
}
func (s *StorageImpl) GetCredentials(userID string) ([]internal.Credential, error) {
	return make([]internal.Credential, 0), nil
}
