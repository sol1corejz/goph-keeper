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
		log.Fatal("failed to create users table:", err)
		return err
	}

	_, err = DBStorage.DB.Exec(`
	CREATE TABLE IF NOT EXISTS credentials (
			uuid UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
			data TEXT NOT NULL,
			meta TEXT
		)
	`)
	if err != nil {
		log.Fatal("failed to create credentials table:", err)
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
func (s *StorageImpl) SaveCredential(cred internal.Credential) error {
	_, err := DBStorage.DB.Exec(`
		INSERT INTO credentials (uuid, user_id, data, meta) VALUES ($1, $2, $3, $4)
	`, cred.ID, cred.UserID, cred.Data, cred.Meta)

	if err != nil {
		log.Info("failed to save credential", err.Error())
		return err
	}

	return nil
}

func (s *StorageImpl) GetCredentials(userID string) ([]internal.Credential, error) {
	rows, err := DBStorage.DB.Query(`
		SELECT uuid, user_id, data, meta FROM credentials WHERE user_id=$1
	`, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("failed to find credentials")
			return nil, ErrNotFound
		}
		log.Info("failed to retrieve credentials", err.Error())

		return nil, err
	}
	defer rows.Close()

	credentials := make([]internal.Credential, 0)
	for rows.Next() {
		var cred internal.Credential
		err := rows.Scan(&cred.ID, &cred.UserID, &cred.Data, &cred.Meta)
		if err != nil {
			log.Info("failed to retrieve credentials", err.Error())
			return nil, err
		}
		credentials = append(credentials, cred)
	}

	if err = rows.Err(); err != nil {
		log.Info("failed to retrieve credentials", err.Error())
		return nil, err
	}

	return credentials, nil
}
