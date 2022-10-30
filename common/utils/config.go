package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Config struct {
	DBDriver            string
	DBSource            string
	ServerAddress       string
	TokenSymmetricKey   string
	AccessTokenDuration time.Duration
}

// These values map to the environment variable names. Should be same as above config file
const (
	DatabaseDriver      = "DB_DRIVER"
	DatabaseURL         = "DATABASE_URL"
	ServerAddress       = "SERVER_ADDRESS"
	TokenSymmetricKey   = "TOKEN_SYMMETRIC_KEY"
	AccessTokenDuration = "ACCESS_TOKEN_DURATION"
)

func LoadConfig() (config Config) {
	config.DBDriver = os.Getenv(DatabaseDriver)
	config.DBSource = os.Getenv(DatabaseURL)
	config.ServerAddress = os.Getenv(ServerAddress)
	config.TokenSymmetricKey = os.Getenv(TokenSymmetricKey)
	accessTokenDurationStr := os.Getenv(AccessTokenDuration)
	var accessTokenDuration time.Duration
	var err error
	if accessTokenDurationStr == "" {
		accessTokenDurationStr = "15m"
	}
	if accessTokenDuration, err = time.ParseDuration(accessTokenDurationStr); err != nil {
		log.Errorln(err)
		accessTokenDuration = time.Minute * 15
	}
	config.AccessTokenDuration = accessTokenDuration
	return config
}
