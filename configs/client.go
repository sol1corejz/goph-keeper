package configs

import (
	"github.com/spf13/viper"
)

type clientConfig struct {
	ServerAddress string `mapstructure:"server_address"`
	SyncInterval  string `mapstructure:"sync_interval"`
	Timeout       string `mapstructure:"timeout"`
}

type clientSecurityConfig struct {
	EncryptionKey string `mapstructure:"encryption_key"`
}

type clientLoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type ClientConfig struct {
	Client   clientConfig         `mapstructure:"client"`
	Security clientSecurityConfig `mapstructure:"security"`
	Logging  clientLoggingConfig  `mapstructure:"logging"`
}

func LoadClientConfig(path string) (*ClientConfig, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config ClientConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
