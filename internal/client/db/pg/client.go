package pg

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/client/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgClient struct {
	masterDBC db.DB
}

// New - конструктор для постгрес клиента.
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к бд: %w", err)
	}

	return &pgClient{masterDBC: &pg{dbc: dbc}}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
