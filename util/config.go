package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env                         string        `mapstructure:"ENVIRONMENT"`
	DBDriver                    string        `mapstructure:"DB_DRIVER"`
	DBSource                    string        `mapstructure:"DB_SOURCE"`
	DBMigrationPath             string        `mapstructure:"DB_MIGRATION_PATH"`
	RediserverAddress           string        `mapstructure:"REDIS_SERVER_ADDRESS"`
	HttpServerAddress           string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress           string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymetricKey            string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	TokenExpiredDuration        time.Duration `mapstructure:"TOKEN_EXPIRED_DURATION"`
	RefreshTokenExpiredDuration time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_DURATION"`
	EmailSenderName             string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress          string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword         string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
