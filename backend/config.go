package backend

import (
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type Config struct {
	DatabaseDriver      string
	DatabaseURL         string
	ServerPort          string
	TokenSymmetricKey   string
	AdminEmail          string
	AdminPassword       string
	AdminName           string
	SendGridAPIKey      string
	TwilioAccountSID    string
	TwilioAuthToken     string
	TwilioPhoneNumber   string
	GinMode             string
	AccessTokenDuration time.Duration
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
	SendGridApiKey      = "SENDGRID_API_KEY"
	TwilioAccountSID    = "TWILIO_ACCOUNT_SID"
	TwilioAuthToken     = "TWILIO_AUTH_TOKEN"
	TwilioPhoneNumber   = "TWILIO_PHONE_NUMBER"
	GinMode             = "GIN_MODE"
)

func LoadConfig() (config Config) {
	var err error
	var accessTokenDuration time.Duration
	accessTokenDurationStr := os.Getenv(AccessTokenDuration)
	if accessTokenDuration, err = time.ParseDuration(accessTokenDurationStr); err != nil {
		accessTokenDuration = time.Hour * 24
	}
	config.DatabaseDriver = os.Getenv(DatabaseDriver)
	config.DatabaseURL = os.Getenv(DatabaseURL)
	config.ServerPort = os.Getenv(ServerPort)
	config.TokenSymmetricKey = os.Getenv(TokenSymmetricKey)
	config.TwilioAccountSID = os.Getenv(TwilioAccountSID)
	config.TwilioPhoneNumber = os.Getenv(TwilioPhoneNumber)
	config.TwilioAuthToken = os.Getenv(TwilioAuthToken)
	config.AdminEmail = os.Getenv(AdminUserEmail)
	config.AdminPassword = os.Getenv(AdminUserPassword)
	config.AdminName = os.Getenv(AdminUserName)
	config.SendGridAPIKey = os.Getenv(SendGridApiKey)
	config.AccessTokenDuration = accessTokenDuration
	config.GinMode = os.Getenv(GinMode)
	gin.SetMode(config.GinMode)
	return config
}
