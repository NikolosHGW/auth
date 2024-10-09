package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env"
)

type grpc struct {
	Host string `env:"GRPC_HOST"`
	Port string `env:"GRPC_PORT"`
}

// NewGRPC - конструктор для переменных для gRPC.
func NewGRPC() (*grpc, error) {
	c := new(grpc)
	err := env.Parse(c)
	if err != nil {
		return nil, fmt.Errorf("не удалось спарсить env: %w", err)
	}

	return c, nil
}

func (c grpc) GetRunAddress() string {
	return net.JoinHostPort(c.Host, c.Port)
}
