// Package configs содержит логику для загрузки и обработки конфигурационных данных,
// используемых в клиентской и серверной частях приложения. Конфигурация загружается из файла,
// с возможностью переопределения значений через переменные окружения.
//
// Пакет включает структуры для настройки сервера, безопасности и логирования, а также
// функцию для загрузки конфигурации из указанного файла.
package configs

import (
	"github.com/spf13/viper"
)

// clientConfig содержит настройки клиента, такие как адрес сервера, интервал синхронизации
// и таймауты для запросов.
type clientConfig struct {
	// ServerAddress — адрес сервера.
	ServerAddress string `mapstructure:"server_address"`

	// SyncInterval — интервал синхронизации клиента с сервером.
	SyncInterval string `mapstructure:"sync_interval"`

	// Timeout — таймаут для запросов.
	Timeout string `mapstructure:"timeout"`
}

// clientSecurityConfig содержит настройки безопасности клиента, такие как ключ для шифрования.
type clientSecurityConfig struct {
	// EncryptionKey — ключ для шифрования.
	EncryptionKey string `mapstructure:"encryption_key"`
}

// clientLoggingConfig содержит настройки логирования клиента, включая уровень логирования
// и путь к файлу логов.
type clientLoggingConfig struct {
	// Level — уровень логирования (например, "debug", "info").
	Level string `mapstructure:"level"`

	// File — путь к файлу, в который будут записываться логи.
	File string `mapstructure:"file"`
}

// ClientConfig объединяет настройки клиента, безопасности и логирования.
type ClientConfig struct {
	// Client — настройки клиента.
	Client clientConfig `mapstructure:"client"`

	// Security — настройки безопасности.
	Security clientSecurityConfig `mapstructure:"security"`

	// Logging — настройки логирования.
	Logging clientLoggingConfig `mapstructure:"logging"`
}

// LoadClientConfig загружает конфигурацию из файла по указанному пути и
// возвращает объект ClientConfig. В случае ошибки возвращает ошибку.
func LoadClientConfig(path string) (*ClientConfig, error) {
	// Устанавливаем файл конфигурации и активируем автоматическое считывание
	// переменных окружения.
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	// Чтение конфигурации из файла.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Разбор конфигурации в структуру ClientConfig.
	var config ClientConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Возвращаем структуру с конфигурацией.
	return &config, nil
}
