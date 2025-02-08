package internal

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2/log"
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

// RegisterGrpc — gRPC-обработчик регистрации пользователя.
func (s *KeeperServer) RegisterGrpc(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Генерация UUID пользователя
	userUuid := uuid.New().String()

	// Хеширование пароля
	hashedPassword, err := HashPassword(in.UserData.Password)
	if err != nil {
		return &pb.RegisterResponse{Error: "Ошибка хеширования пароля"}, err
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
		return &pb.RegisterResponse{Error: "Ошибка сохранения в БД"}, err
	}

	// Генерация токена
	token, err := auth.GenerateToken(s.Config, userUuid)
	if err != nil {
		return &pb.RegisterResponse{Error: "Ошибка генерации токена"}, err
	}

	// Возвращаем успешный ответ
	return &pb.RegisterResponse{Token: token}, nil
}

// LoginGrpc — gRPC-обработчик авторизации пользователя.
func (s *KeeperServer) LoginGrpc(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {

	// Входные данные для авторизации
	loginData := models.AuthPayload{
		Username: in.UserData.Username,
		Password: in.UserData.Password,
	}

	// Получение пользователя из базы данных
	userData, err := storage.DBStorage.GetUser(loginData.Username)
	if err != nil {
		if errors.Is(storage.ErrNotFound, err) {
			return &pb.LoginResponse{
				Error: "Неправильный логин или пароль",
			}, err
		}
		return &pb.LoginResponse{
			Error: err.Error(),
		}, err
	}

	// Сравнение пароля из входных данных с паролем из базы данных
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password))
	if err != nil {
		return &pb.LoginResponse{
			Error: "Неправильный логин или пароль",
		}, err
	}

	// Генерация токена аутентификации
	token, err := auth.GenerateToken(s.Config, userData.ID)

	// Отправка успешного ответа
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *KeeperServer) AddCredentialsGrpc(ctx context.Context, in *pb.AddCredentialsRequest) (*pb.AddCredentialsResponse, error) {
	// Получение токена
	token := in.Token
	if token == "" {
		return &pb.AddCredentialsResponse{
			Error: "Неавторизован",
		}, errors.New("unauthorized")
	}

	// Парсинг входных данных
	credentialsPayload := models.CredentialPayload{
		Data: in.Credentials.Data,
		Meta: in.Credentials.Meta,
	}

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(s.Config, token)
	if err != nil {
		log.Info("token is invalid")
		return &pb.AddCredentialsResponse{
			Error: "Не валидный токен аутентификации",
		}, errors.New("invalid token")
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
		return &pb.AddCredentialsResponse{
			Error: "Ошибка добавления данных",
		}, errors.New("failed to save credential data")
	}

	// Отправка успешного ответа
	return &pb.AddCredentialsResponse{}, nil
}

func (s *KeeperServer) EditCredentialsGrpc(ctx context.Context, in *pb.EditCredentialsRequest) (*pb.EditCredentialsResponse, error) {
	// Получение токена
	token := in.Token
	if token == "" {
		return &pb.EditCredentialsResponse{
			Error: "Неавторизован",
		}, errors.New("unauthorized")
	}

	// Парсинг входных данных
	credentialsPayload := models.CredentialPayload{
		Data: in.Credentials.Data,
		Meta: in.Credentials.Meta,
	}

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(s.Config, token)
	if err != nil {
		log.Info("token is invalid")
		return &pb.EditCredentialsResponse{
			Error: "Не валидный токен аутентификации",
		}, errors.New("invalid token")
	}

	// Подготовка данных для сохранения в базе данных
	credentialsData := models.Credential{
		ID:     uuid.New().String(),
		UserID: userID,
		Data:   credentialsPayload.Data,
		Meta:   credentialsPayload.Meta,
	}

	// Сохранение учетных данных в базе данных
	err = storage.DBStorage.EditCredential(credentialsData)
	if err != nil {
		return &pb.EditCredentialsResponse{
			Error: "Ошибка обновления данных",
		}, errors.New("failed to edit credential data")
	}

	// Отправка успешного ответа
	return &pb.EditCredentialsResponse{}, nil

}

func (s *KeeperServer) GetCredentialsGrpc(ctx context.Context, in *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	// Получение токена
	token := in.Token
	if token == "" {
		return &pb.GetCredentialsResponse{
			Error: "Неавторизован",
		}, errors.New("unauthorized")
	}

	// Проверка авторизации
	userID, err := auth.CheckIsAuthorized(s.Config, token)
	if err != nil {
		log.Info("token is invalid")
		return &pb.GetCredentialsResponse{
			Error: "Не валидный токен аутентификации",
		}, errors.New("invalid token")
	}

	// Получение учетных данных пользователя из базы данных
	var credentialsData []models.Credential
	credentialsData, err = storage.DBStorage.GetCredentials(userID)
	if err != nil {
		return &pb.GetCredentialsResponse{
			Error: "Ошибка получения данных",
		}, errors.New("failed to retrieve credentials")
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
	return &pb.GetCredentialsResponse{
		Credentials: credentials,
	}, nil
}
