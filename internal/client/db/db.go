package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler func(ctx context.Context) error

// Client - клиента для работы с БД.
type Client interface {
	DB() DB
	Close() error
}

// Query - обёртка над запросом, хранящая имя запроса и сам запрос.
type Query struct {
	Name     string
	QueryRaw string
}

// SQLExecer комбинирует NamedExecer и QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer - интерфейс для работы с именованными запросами с помощью тегов в структурах.
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest any, q Query, args ...any) error
	ScanAllContext(ctx context.Context, dest any, q Query, args ...any) error
}

// QueryExecer - интерфейс для работы с обычными запросами.
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...any) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...any) pgx.Row
}

// Pinger - интерфейс для проверки соединения с БД.
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB - интерфейс для работы с БД.
type DB interface {
	SQLExecer
	Pinger
	Transactor
	Close()
}

// Transactor - интерфейс для работы с транзакциями.
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}
