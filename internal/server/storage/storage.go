package internal

import (
	"database/sql"
	"errors"
	log "github.com/gofiber/fiber/v2/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sol1corejz/goph-keeper/configs"
	internal "github.com/sol1corejz/goph-keeper/internal/common/models"
)

// Storage - Интерфейс для абстракции хранилища.
// Это позволяет легко расширять функциональность и использовать разные виды хранения данных.
// Также интерфейс упрощает написание тестов, предоставляя возможность создания мокового хранилища.
type Storage interface {
	// CreateUser создаёт пользователя в хранилище.
	CreateUser(user internal.User) error
	// GetUser возвращает данные пользователя по имени пользователя.
	GetUser(username string) (internal.User, error)
	// SaveCredential сохраняет учётные данные для пользователя.
	SaveCredential(userID string, cred internal.Credential) error
	// GetCredentials возвращает список учётных данных пользователя по его ID.
	GetCredentials(userID string) ([]internal.Credential, error)
}

// StorageImpl - Реализация интерфейса Storage с использованием базы данных.
type StorageImpl struct {
	DB *sql.DB
}

// DBStorage - Глобальная переменная для доступа к методам хранилища.
var DBStorage StorageImpl

// ErrNotFound - Ошибка, возвращаемая, если данные не найдены.
var ErrNotFound = errors.New("not found")

// ConnectDB подключается к базе данных и создаёт необходимые таблицы.
// Возвращает ошибку, если подключение или создание таблиц не удалось.
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

	// Создание таблицы пользователей
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

	// Создание таблицы учётных данных
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

// CreateUser создаёт нового пользователя в базе данных.
// Принимает структуру пользователя, содержащую ID, имя пользователя и пароль.
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

// GetUser возвращает данные пользователя по его имени пользователя.
// Если пользователь не найден, возвращается ошибка ErrNotFound.
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

// SaveCredential сохраняет учётные данные пользователя в базу данных.
// Принимает структуру Credential, содержащую ID, ID пользователя, данные и метаинформацию.
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

// GetCredentials возвращает список учётных данных пользователя по его ID.
// Если данные не найдены, возвращается ошибка ErrNotFound.
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
