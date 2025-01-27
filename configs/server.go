package configs

import (
	"github.com/spf13/viper"
)

// serverConfig описывает настройки сервера, включая адрес и таймауты.
type serverConfig struct {
	Address      string `mapstructure:"address"`       // Адрес сервера (например, ":8080").
	ReadTimeout  string `mapstructure:"read_timeout"`  // Таймаут на чтение запроса.
	WriteTimeout string `mapstructure:"write_timeout"` // Таймаут на запись ответа.
	IdleTimeout  string `mapstructure:"idle_timeout"`  // Таймаут бездействия соединения.
}

// serverStorageConfig описывает настройки хранилища данных.
type serverStorageConfig struct {
	Type             string `mapstructure:"type"`              // Тип хранилища (например, "postgres", "file").
	ConnectionString string `mapstructure:"connection_string"` // Строка подключения к базе данных.
	FilePath         string `mapstructure:"file_path"`         // Путь к файлу (используется для файлового хранилища).
}

// serverSecurityConfig описывает параметры безопасности сервера.
type serverSecurityConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`     // Секрет для подписи JWT.
	EncryptionKey string `mapstructure:"encryption_key"` // Ключ для шифрования данных.
}

// serverLoggingConfig описывает параметры логирования.
type serverLoggingConfig struct {
	Level string `mapstructure:"level"` // Уровень логирования (например, "debug", "info", "error").
	File  string `mapstructure:"file"`  // Путь к файлу для записи логов.
}

// ServerConfig объединяет все настройки сервера.
type ServerConfig struct {
	Server   serverConfig         `mapstructure:"server"`   // Настройки сервера.
	Storage  serverStorageConfig  `mapstructure:"storage"`  // Настройки хранилища.
	Security serverSecurityConfig `mapstructure:"security"` // Настройки безопасности.
	Logging  serverLoggingConfig  `mapstructure:"logging"`  // Настройки логирования.
}

// LoadServerConfig загружает конфигурацию сервера из указанного файла.
//
// path - путь к файлу конфигурации.
//
// Функция использует библиотеку viper для загрузки настроек из файла и окружения.
// Она возвращает указатель на объект ServerConfig или ошибку, если загрузка не удалась.
func LoadServerConfig(path string) (*ServerConfig, error) {
	viper.SetConfigFile(path) // Установка пути к файлу конфигурации.
	viper.AutomaticEnv()      // Автоматическое чтение переменных окружения.

	// Чтение конфигурационного файла.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Декодирование конфигурации в структуру ServerConfig.
	var config ServerConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
