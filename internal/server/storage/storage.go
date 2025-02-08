package internal

import (
	"database/sql"
	"errors"
	log "github.com/gofiber/fiber/v2/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sol1corejz/goph-keeper/configs"
	internal "github.com/sol1corejz/goph-keeper/internal/server/models"
)

// Storage - интерфейс хранилища данных, предоставляющий методы для работы с пользователями и учетными данными.
// Используется для расширяемости и удобства тестирования.
type Storage interface {
	// CreateUser добавляет нового пользователя в хранилище.
	CreateUser(user internal.User) error
	// GetUser получает данные пользователя по имени пользователя.
	GetUser(username string) (internal.User, error)
	// SaveCredential сохраняет учетные данные пользователя.
	SaveCredential(cred internal.Credential) error
	// EditCredential сохраняет учетные данные пользователя.
	EditCredential(cred internal.Credential) error
	// GetCredentials возвращает все учетные данные пользователя.
	GetCredentials(userID string) ([]internal.Credential, error)
}

// StorageImpl - реализация интерфейса Storage, использующая базу данных PostgreSQL.
type StorageImpl struct {
	DB *sql.DB
}

// DBStorage - глобальный объект для работы с базой данных.
var DBStorage StorageImpl

// ErrNotFound - ошибка, возвращаемая при отсутствии данных.
var ErrNotFound = errors.New("not found")

// ErrAlreadyExists - ошибка, возвращаемая при существовании данных.
var ErrAlreadyExists = errors.New("already exists")

// ConnectDB устанавливает соединение с базой данных и создает необходимые таблицы.
func ConnectDB(cfg *configs.ServerConfig) error {
	if cfg.Storage.ConnectionString == "" {
		return errors.New("no connection string provided")
	}

	// Открываем соединение с базой данных PostgreSQL
	db, err := sql.Open("pgx", cfg.Storage.ConnectionString)
	if err != nil {
		log.Fatal(err)
		return err
	}

	DBStorage.DB = db

	// Создаем таблицу пользователей, если она отсутствует
	_, err = DBStorage.DB.Exec(`
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

	// Создаем таблицу учетных данных пользователей
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

// CreateUser добавляет нового пользователя в базу данных.
func (s *StorageImpl) CreateUser(user internal.User) error {
	var exists bool
	err := DBStorage.DB.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)
    `, user.Username).Scan(&exists)
	if err != nil {
		log.Error("failed to check user: ", err)
		return err
	}
	if exists {
		return ErrAlreadyExists
	}

	_, err = DBStorage.DB.Exec(`
		INSERT INTO users (uuid, username, password) VALUES ($1, $2, $3)
	`, user.ID, user.Username, user.Password)

	if err != nil {
		log.Info("failed to create user")
		return err
	}

	return nil
}

// GetUser получает данные пользователя по имени пользователя.
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

// SaveCredential сохраняет учетные данные пользователя в базе данных.
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

// EditCredential обновляет учетные данные пользователя в базе данных.
func (s *StorageImpl) EditCredential(cred internal.Credential) error {
	_, err := DBStorage.DB.Exec(`
		UPDATE credentials SET data = $1, meta = $2 WHERE uuid = $3
	`, cred.Data, cred.Meta, cred.ID)

	if err != nil {
		log.Info("failed to save credential", err.Error())
		return err
	}

	return nil
}

// GetCredentials получает все учетные данные, принадлежащие пользователю.
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
