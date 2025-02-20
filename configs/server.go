package configs

import (
	"github.com/spf13/viper"
)

// serverConfig содержит настройки сервера, такие как адрес, таймауты на чтение,
// запись и таймаут простоя.
type serverConfig struct {
	// Address — адрес, на котором работает сервер.
	Address string `mapstructure:"address"`

	// ReadTimeout — таймаут для чтения данных от клиента.
	ReadTimeout string `mapstructure:"read_timeout"`

	// WriteTimeout — таймаут для записи данных на клиента.
	WriteTimeout string `mapstructure:"write_timeout"`

	// IdleTimeout — таймаут простоя, после которого соединение с клиентом закрывается.
	IdleTimeout string `mapstructure:"idle_timeout"`
}

// serverStorageConfig содержит настройки хранилища данных для сервера,
// такие как тип хранилища, строка подключения и путь к файлу для хранения данных.
type serverStorageConfig struct {
	// Type — тип хранилища (например, "postgres", "mysql", "file").
	Type string `mapstructure:"type"`

	// ConnectionString — строка подключения к хранилищу данных.
	ConnectionString string `mapstructure:"connection_string"`

	// FilePath — путь к файлу для хранения данных, если используется файловое хранилище.
	FilePath string `mapstructure:"file_path"`
}

// serverSecurityConfig содержит настройки безопасности для сервера,
// такие как секретный ключ для JWT и ключ для шифрования.
type serverSecurityConfig struct {
	// JWTSecret — секретный ключ для подписания JWT.
	JWTSecret string `mapstructure:"jwt_secret"`

	// EncryptionKey — ключ для шифрования данных.
	EncryptionKey string `mapstructure:"encryption_key"`
}

// serverLoggingConfig содержит настройки логирования для сервера,
// включая уровень логирования и путь к файлу для записи логов.
type serverLoggingConfig struct {
	// Level — уровень логирования (например, "debug", "info").
	Level string `mapstructure:"level"`

	// File — путь к файлу, в который будут записываться логи.
	File string `mapstructure:"file"`
}

// ServerConfig объединяет все настройки сервера, хранилища, безопасности и логирования.
type ServerConfig struct {
	// Server — настройки сервера.
	Server serverConfig `mapstructure:"server"`

	// Storage — настройки хранилища данных.
	Storage serverStorageConfig `mapstructure:"storage"`

	// Security — настройки безопасности.
	Security serverSecurityConfig `mapstructure:"security"`

	// Logging — настройки логирования.
	Logging serverLoggingConfig `mapstructure:"logging"`
}

// LoadServerConfig загружает конфигурацию из файла по указанному пути и
// возвращает объект ServerConfig. В случае ошибки возвращает ошибку.
func LoadServerConfig(path string) (*ServerConfig, error) {
	// Устанавливаем файл конфигурации и активируем автоматическое считывание
	// переменных окружения.
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	// Чтение конфигурации из файла.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Разбор конфигурации в структуру ServerConfig.
	var config ServerConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Возвращаем структуру с конфигурацией.
	return &config, nil
}
