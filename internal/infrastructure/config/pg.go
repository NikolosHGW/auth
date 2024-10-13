package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type pg struct {
	DatabaseDSN string `env:"DATABASE_DSN"`
}

// NewPG - конструктор переменных для постгреса.
func NewPG() (*pg, error) {
	c := new(pg)
	err := env.Parse(c)
	if err != nil {
		return nil, fmt.Errorf("не удалось спарсить env: %w", err)
	}

	return c, nil
}

func (c pg) GetDatabaseDSN() string {
	return c.DatabaseDSN
}
