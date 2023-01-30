package backend

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Config struct {
	DatabaseDriver      string
	DatabaseURL         string
	ServerPort          string
	TokenSymmetricKey   string
	AccessTokenDuration time.Duration
	AdminEmail          string
	AdminPassword       string
	AdminName           string
}

// These values map to the environment variable names. Should be same as above config file
const (
	DatabaseDriver      = "DB_DRIVER"
	DatabaseURL         = "DATABASE_URL"
	ServerPort          = "PORT"
	TokenSymmetricKey   = "TOKEN_SYMMETRIC_KEY"
	AccessTokenDuration = "ACCESS_TOKEN_DURATION"
	AdminUserEmail      = "ADMIN_USER_EMAIL"
	AdminUserPassword   = "ADMIN_USER_PASSWORD"
	AdminUserName       = "ADMIN_USER_NAME"
)

const DefaultDatabaseDriver = "mongodb"
const DefaultDatabaseURL = "mongodb://root:example@localhost:27017/?"
const DefaultServerPort = "8080"
const DefaultTokenSymmetricKey = "12345678901234567890123456789012"
const DefaultAccessTokenDuration = "24h"

func LoadConfig() (config Config) {
	if config.DatabaseDriver = os.Getenv(DatabaseDriver); config.DatabaseDriver == "" {
		config.DatabaseDriver = DefaultDatabaseDriver
	}
	if config.DatabaseURL = os.Getenv(DatabaseURL); config.DatabaseURL == "" {
		config.DatabaseURL = DefaultDatabaseURL
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
	config.AdminEmail = os.Getenv(AdminUserEmail)
	config.AdminPassword = os.Getenv(AdminUserPassword)
	config.AdminName = os.Getenv(AdminUserName)
	config.AccessTokenDuration = accessTokenDuration
	return config
}
