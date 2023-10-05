package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	HttpAddress          string        `mapstructure:"http_address"`
	TokenSymmetricKey    string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	Storage              StorageConfig `mapstructure:"storage"`
}
type StorageConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port" `
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
