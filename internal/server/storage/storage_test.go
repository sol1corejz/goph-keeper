package internal_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	models "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestCreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	store := &storage.StorageImpl{DB: mockDB}
	fmt.Println("store.DB:", store.DB)
	if store.DB == nil {
		t.Fatal("database connection is nil")
	}

	user := models.User{
		ID:       uuid.New().String(),
		Username: "testuser",
		Password: "password123",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)")).
		WithArgs(user.Username).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (uuid, username, password) VALUES ($1, $2, $3)")).
		WithArgs(user.ID, user.Username, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	store := &storage.StorageImpl{DB: mockDB}
	user := models.User{
		ID:       uuid.New().String(),
		Username: "testuser",
		Password: "password123",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username=$1")).
		WithArgs(user.Username).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "username", "password"}).
			AddRow(user.ID, user.Username, user.Password))

	result, err := store.GetUser(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveCredential(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	store := &storage.StorageImpl{DB: mockDB}
	cred := models.Credential{
		ID:     uuid.New().String(),
		UserID: uuid.New().String(),
		Data:   "secure data",
		Meta:   "metadata",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO credentials")).
		WithArgs(cred.ID, cred.UserID, cred.Data, cred.Meta).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.SaveCredential(cred)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEditCredential(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	store := &storage.StorageImpl{DB: mockDB}
	cred := models.Credential{
		ID:   uuid.New().String(),
		Data: "updated data",
		Meta: "updated meta",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE credentials SET data = $1, meta = $2 WHERE uuid = $3")).
		WithArgs(cred.Data, cred.Meta, cred.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.EditCredential(cred)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCredentials(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	store := &storage.StorageImpl{DB: mockDB}
	userID := uuid.New().String()
	cred1 := models.Credential{ID: uuid.New().String(), UserID: userID, Data: "data1", Meta: "meta1"}
	cred2 := models.Credential{ID: uuid.New().String(), UserID: userID, Data: "data2", Meta: "meta2"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT uuid, user_id, data, meta FROM credentials WHERE user_id=$1")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "user_id", "data", "meta"}).
			AddRow(cred1.ID, cred1.UserID, cred1.Data, cred1.Meta).
			AddRow(cred2.ID, cred2.UserID, cred2.Data, cred2.Meta))

	result, err := store.GetCredentials(userID)
	assert.NoError(t, err)
	assert.Equal(t, []models.Credential{cred1, cred2}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
