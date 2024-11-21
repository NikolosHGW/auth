package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Load - подгружает .env файл из корня.
func Load() error {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./.env"
	}

	err := godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("не удалось загрузить файл .env: %w", err)
	}

	return nil
}

// GRPCConfiger - контракт для конфига gRPC.
type GRPCConfiger interface {
	GetRunAddress() string
}

// PGConfiger - контракт для конфига постгреса.
type PGConfiger interface {
	GetDatabaseDSN() string
}

// HTTPConfiger - контракт для конфига http-сервера.
type HTTPConfiger interface {
	HTTPRunAddress() string
}
