package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver              string        `mapstructure:"DB_DRIVER"`
	DBSource              string        `mapstructure:"DB_SOURCE"`
	MigrationURL          string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress     string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress     string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	HTTPAuthServerAddress string        `mapstructure:"HTTP_AUTH_SERVER_ADDRESS"`
	GRPCAuthServerAddress string        `mapstructure:"GRPC_AUTH_SERVER_ADDRESS"`
	TokenSymmetricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	OauthClientID         string        `mapstructure:"OAUTH_CLIENT_ID"`
	OauthClientSecret     string        `mapstructure:"OAUTH_CLIENT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, yaml

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
