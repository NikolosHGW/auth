package config

import (
	"fmt"
	"log"
	"net"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	GRPCHost    string `env:"GRPC_HOST"`
	GRPCPort    string `env:"GRPC_PORT"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func (c *config) initEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("не удалось загрузить файл .env: %w", err)
	}

	err = env.Parse(c)
	if err != nil {
		return fmt.Errorf("не удалось спарсить env: %w", err)
	}

	return nil
}

// NewConfig - создание конфига. Использовать позже.
func NewConfig() *config {
	cfg := new(config)

	if err := cfg.initEnv(); err != nil {
		log.Fatalf("Ошибка при инициализации переменных окружения: %v", err)
	}

	return cfg
}

func (c config) GetRunAddress() string {
	return net.JoinHostPort(c.GRPCHost, c.GRPCPort)
}

func (c config) GetDatabaseDSN() string {
	return c.DatabaseDSN
}
