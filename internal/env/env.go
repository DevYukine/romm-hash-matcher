package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"romm-hash-matcher/internal/logging"
)

type config struct {
	RommUsername string `env:"ROMM_USERNAME,required"`
	RommPassword string `env:"ROMM_PASSWORD,required"`
	RommUrl      string `env:"ROMM_URL,required"`
}

var Config config

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		// If .env file is not found, we can still proceed with environment variables
	}

	err = env.Parse(&Config)

	if err != nil {
		logging.Logger.Fatal("Error parsing environment variables", zap.Error(err))
	}
}
