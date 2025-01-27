package configs

import (
	"github.com/spf13/viper"
)

// clientConfig содержит настройки клиента, такие как адрес сервера, интервал синхронизации и таймаут.
type clientConfig struct {
	ServerAddress string `mapstructure:"server_address"` // Адрес сервера.
	SyncInterval  string `mapstructure:"sync_interval"`  // Интервал синхронизации.
	Timeout       string `mapstructure:"timeout"`        // Таймаут.
}

// clientSecurityConfig содержит настройки безопасности клиента, включая ключ шифрования.
type clientSecurityConfig struct {
	EncryptionKey string `mapstructure:"encryption_key"` // Ключ шифрования.
}

// clientLoggingConfig содержит настройки логирования клиента, включая уровень логирования и файл логов.
type clientLoggingConfig struct {
	Level string `mapstructure:"level"` // Уровень логирования.
	File  string `mapstructure:"file"`  // Путь к файлу логов.
}

// ClientConfig объединяет настройки клиента, безопасности и логирования.
type ClientConfig struct {
	Client   clientConfig         `mapstructure:"client"`   // Настройки клиента.
	Security clientSecurityConfig `mapstructure:"security"` // Настройки безопасности.
	Logging  clientLoggingConfig  `mapstructure:"logging"`  // Настройки логирования.
}

// LoadClientConfig загружает конфигурацию клиента из указанного файла.
// Он автоматически считывает переменные окружения и пытается загрузить настройки из конфигурационного файла.
// Возвращает указатель на структуру ClientConfig или ошибку в случае неудачи.
func LoadClientConfig(path string) (*ClientConfig, error) {
	// Устанавливаем путь к конфигурационному файлу.
	viper.SetConfigFile(path)
	// Автоматически загружаем переменные окружения.
	viper.AutomaticEnv()

	// Читаем конфигурацию из файла.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config ClientConfig
	// Десериализуем данные из конфигурации в структуру.
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Возвращаем загруженную конфигурацию.
	return &config, nil
}
