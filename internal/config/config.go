package config

import (
	"os"

	"github.com/dveselov/juno/internal/authentication/providers"
)

type AppConfig struct {
	Port      string
	Secret    string
	Providers map[string]providers.AuthenticationByCodeProvider
}

func GetAppConfig() AppConfig {
	return AppConfig{
		Port:   os.Getenv("APP_PORT"),
		Secret: os.Getenv("APP_SECRET"),
		Providers: map[string]providers.AuthenticationByCodeProvider{
			"call": providers.UCallerProvider{
				ServiceID:  os.Getenv("UCALLER_SERVICE_ID"),
				SecretKey:  os.Getenv("UCALLER_SECRET_KEY"),
				BaseAPIUrl: "https://api.ucaller.ru/v1.0/initCall",
			},
		},
	}
}
