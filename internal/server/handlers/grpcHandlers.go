package internal

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sol1corejz/goph-keeper/configs"
	"github.com/sol1corejz/goph-keeper/internal/server/auth"
	models "github.com/sol1corejz/goph-keeper/internal/server/models"
	storage "github.com/sol1corejz/goph-keeper/internal/server/storage"
	pb "github.com/sol1corejz/goph-keeper/proto"
	"golang.org/x/crypto/bcrypt"
)

// KeeperServer реализует gRPC Keeper.
type KeeperServer struct {
	pb.UnimplementedKeeperServer
	Config *configs.ServerConfig
}

// Register — gRPC-обработчик регистрации пользователя.
func (s *KeeperServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Генерация UUID пользователя
	userUuid := uuid.New().String()

	// Создание ответа с пустым значением, чтобы он всегда был возвращен
	resp := &pb.RegisterResponse{}

	// Хеширование пароля
	hashedPassword, err := HashPassword(in.UserData.Password)
	if err != nil {
		resp.Error = "Ошибка хеширования пароля"
		return resp, err
	}

	// Создание пользователя
	userData := models.User{
		ID:       userUuid,
		Username: in.UserData.Username,
		Password: hashedPassword,
	}

	// Сохранение в БД
	err = storage.DBStorage.CreateUser(userData)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			resp.Error = "Пользователь уже зарегистрирован"
			return resp, err
		}
		resp.Error = "Ошибка сохранения в БД"
		return resp, err
	}

	// Генерация токена
	token, err := auth.GenerateToken(s.Config, userUuid)
	if err != nil {
		resp.Error = "Ошибка генерации токена"
		return resp, err
	}

	// Возвращаем успешный ответ
	resp.Token = token
	return resp, nil
}

// Login — gRPC-обработчик авторизации пользователя.
func (s *KeeperServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Создание ответа с пустым значением, чтобы он всегда был возвращен
	resp := &pb.LoginResponse{}

	// Входные данные для авторизации
	loginData := models.AuthPayload{
		Username: in.UserData.Username,
		Password: in.UserData.Password,
	}

	// Получение пользователя из базы данных
	userData, err := storage.DBStorage.GetUser(loginData.Username)
	if err != nil {
		if errors.Is(storage.ErrNotFound, err) {
			resp.Error = "Неправильный логин или пароль"
			return resp, err
		}
		resp.Error = err.Error()
		return resp, err
	}

	// Сравнение пароля из входных данных с паролем из базы данных
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password))
	if err != nil {
		resp.Error = "Неправильный логин или пароль"
		return resp, err
	}

	// Генерация токена аутентификации
	token, err := auth.GenerateToken(s.Config, userData.ID)

	// Отправка успешного ответа
	resp.Token = token
	return resp, nil
}

// AddCredentials — gRPC-обработчик для добавления данных пользователя.
func (s *KeeperServer) AddCredentials(ctx context.Context, in *pb.AddCredentialsRequest) (*pb.AddCredentialsResponse, error) {
	// Создание ответа с пустым значением, чтобы он всегда был возвращен
	resp := &pb.AddCredentialsResponse{}

	// Получение токена
	token := in.Token
	if token == "" {
		resp.Error = "Неавторизован"
		return resp, errors.New("unauthorized")
	}

	// Парсинг входных данных
	credentialsPayload := models.CredentialPayload{
		Data: in.Credentials.Data,
		Meta: in.Credentials.Meta,
	}

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(s.Config, token)
	if err != nil {
		resp.Error = "Не валидный токен аутентификации"
		return resp, errors.New("invalid token")
	}

	// Подготовка данных для сохранения в базе данных
	credentialsData := models.Credential{
		ID:     uuid.New().String(),
		UserID: userID,
		Data:   credentialsPayload.Data,
		Meta:   credentialsPayload.Meta,
	}

	// Сохранение учетных данных в базе данных
	err = storage.DBStorage.SaveCredential(credentialsData)
	if err != nil {
		resp.Error = "Ошибка добавления данных"
		return resp, errors.New("failed to save credential data")
	}

	// Отправка успешного ответа
	return resp, nil
}

// EditCredentials — gRPC-обработчик для редактирования данных пользователя.
func (s *KeeperServer) EditCredentials(ctx context.Context, in *pb.EditCredentialsRequest) (*pb.EditCredentialsResponse, error) {
	// Создание ответа с пустым значением, чтобы он всегда был возвращен
	resp := &pb.EditCredentialsResponse{}

	// Получение токена
	token := in.Token
	if token == "" {
		resp.Error = "Неавторизован"
		return resp, errors.New("unauthorized")
	}

	// Парсинг входных данных
	credentialsPayload := models.CredentialPayload{
		Data: in.Credentials.Data,
		Meta: in.Credentials.Meta,
	}

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(s.Config, token)
	if err != nil {
		resp.Error = "Не валидный токен аутентификации"
		return resp, errors.New("invalid token")
	}

	// Подготовка данных для сохранения в базе данных
	credentialsData := models.Credential{
		ID:     in.Id,
		UserID: userID,
		Data:   credentialsPayload.Data,
		Meta:   credentialsPayload.Meta,
	}

	// Сохранение учетных данных в базе данных
	err = storage.DBStorage.EditCredential(credentialsData)
	if err != nil {
		resp.Error = "Ошибка обновления данных"
		return resp, errors.New("failed to edit credential data")
	}

	// Отправка успешного ответа
	return resp, nil

}

// GetCredentials — gRPC-обработчик для получения данных польхователя.
func (s *KeeperServer) GetCredentials(ctx context.Context, in *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	// Создание ответа с пустым значением, чтобы он всегда был возвращен
	resp := &pb.GetCredentialsResponse{}

	// Получение токена
	token := in.Token
	if token == "" {
		resp.Error = "Неавторизован"
		return resp, errors.New("unauthorized")
	}

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(s.Config, token)
	if err != nil {
		resp.Error = "Не валидный токен аутентификации"
		return resp, errors.New("invalid token")
	}

	// Получение учетных данных пользователя из базы данных
	var credentialsData []models.Credential
	credentialsData, err = storage.DBStorage.GetCredentials(userID)
	if err != nil {
		resp.Error = "Ошибка получения данных"
		return resp, errors.New("failed to retrieve credentials")
	}

	// Преобразование данных для отправки ответа
	credentials := make([]*pb.Credentials, 0)
	for _, credential := range credentialsData {
		credentials = append(credentials, &pb.Credentials{
			Data: credential.Data,
			Meta: credential.Meta,
		})
	}

	// Отправка учетных данных в ответе
	resp.Credentials = credentials
	return resp, nil
}
