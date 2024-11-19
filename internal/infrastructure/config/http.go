package config

import (
	"fmt"
	"net"

	"github.com/caarlos0/env"
)

type http struct {
	Host string `env:"HTTP_HOST"`
	Port string `env:"HTTP_PORT"`
}

// NewHTTP - конструктор для переменных для HTTP-сервера.
func NewHTTP() (*http, error) {
	c := new(http)
	err := env.Parse(c)
	if err != nil {
		return nil, fmt.Errorf("не удалось спарсить env: %w", err)
	}

	return c, nil
}

func (c http) HTTPRunAddress() string {
	return net.JoinHostPort(c.Host, c.Port)
}
