package sqlc

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Config struct {
	DatabaseDriver      string
	DatabaseURL         string
	ServerHost          string
	ServerPort          string
	TokenSymmetricKey   string
	AccessTokenDuration time.Duration
}

const (
	DatabaseDriver      = "DB_DRIVER"
	DatabaseURL         = "DATABASE_URL"
	ServerAddress       = "HOST"
	ServerPort          = "PORT"
	TokenSymmetricKey   = "TOKEN_SYMMETRIC_KEY"
	AccessTokenDuration = "ACCESS_TOKEN_DURATION"
)

const DefaultDatabaseDriver = "postgres"
const DefaultDatabaseURL = "postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable"
const DefaultServerHost = "0.0.0.0"
const DefaultServerPort = "80"
const DefaultTokenSymmetricKey = "12345678901234567890123456789012"
const DefaultAccessTokenDuration = "24h"

func LoadConfig() (config Config) {
	if config.DatabaseDriver = os.Getenv(DatabaseDriver); config.DatabaseDriver == "" {
		config.DatabaseDriver = DefaultDatabaseDriver
	}
	if config.DatabaseURL = os.Getenv(DatabaseURL); config.DatabaseURL == "" {
		config.DatabaseURL = DefaultDatabaseURL
	}
	if config.ServerHost = os.Getenv(ServerAddress); config.ServerHost == "" {
		config.ServerHost = DefaultServerHost
	}
	if config.ServerPort = os.Getenv(ServerPort); config.ServerPort == "" {
		config.ServerPort = DefaultServerPort
	}
	if config.TokenSymmetricKey = os.Getenv(TokenSymmetricKey); config.TokenSymmetricKey == "" {
		config.TokenSymmetricKey = DefaultTokenSymmetricKey
	}
	accessTokenDurationStr := os.Getenv(AccessTokenDuration)
	var accessTokenDuration time.Duration
	var err error
	if accessTokenDurationStr == "" {
		accessTokenDurationStr = DefaultAccessTokenDuration
	}
	if accessTokenDuration, err = time.ParseDuration(accessTokenDurationStr); err != nil {
		log.Errorln(err)
		accessTokenDuration = time.Minute * 15
	}

	config.AccessTokenDuration = accessTokenDuration
	return config
}
