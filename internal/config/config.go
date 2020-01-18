package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"

	"github.com/dveselov/juno/internal/authentication/providers"
)

type AppContext struct {
	Config AppConfig
	echo.Context
}

func (c *AppContext) GetDB() *sqlx.DB {
	return c.Config.Database
}

func (c *AppContext) GetAuthenticationProvider() providers.AuthenticationByCodeProvider {
	return c.Config.Provider
}

type AppConfig struct {
	Database               *sqlx.DB
	Port                   string
	ListenAddr             string
	SigningKey             string
	AuthenticationTokenTTL string
	Provider               providers.AuthenticationByCodeProvider
	VerificationCodeTTL    string
}

func GetAppConfig(db *sqlx.DB) AppConfig {
	provider := providers.UCallerProvider{
		ServiceID:  os.Getenv("UCALLER_SERVICE_ID"),
		SecretKey:  os.Getenv("UCALLER_SECRET_KEY"),
		BaseAPIUrl: "https://api.ucaller.ru/v1.0",
	}
	provider.Database = db
	return AppConfig{
		Database:               db,
		ListenAddr:             fmt.Sprintf(":%s", os.Getenv("APP_LISTEN_PORT")),
		SigningKey:             os.Getenv("APP_SIGNING_KEY"),
		AuthenticationTokenTTL: os.Getenv("APP_AUTHENTICATION_TOKEN_TTL"),
		Provider:               provider,
		VerificationCodeTTL:    "10 minutes",
	}
}
