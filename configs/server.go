package configs

import (
	"github.com/spf13/viper"
)

type serverConfig struct {
	Address      string `mapstructure:"address"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	IdleTimeout  string `mapstructure:"idle_timeout"`
}

type serverStorageConfig struct {
	Type             string `mapstructure:"type"`
	ConnectionString string `mapstructure:"connection_string"`
	FilePath         string `mapstructure:"file_path"`
}

type serverSecurityConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	EncryptionKey string `mapstructure:"encryption_key"`
}

type serverLoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type ServerConfig struct {
	Server   serverConfig         `mapstructure:"server"`
	Storage  serverStorageConfig  `mapstructure:"storage"`
	Security serverSecurityConfig `mapstructure:"security"`
	Logging  serverLoggingConfig  `mapstructure:"logging"`
}

func LoadServerConfig(path string) (*ServerConfig, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
