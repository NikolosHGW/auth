package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/NikolosHGW/auth/internal/client/db"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pg struct {
	dbc *pgxpool.Pool
}

// NewPG - конструктор реализации постгрес клиента.
func NewPG(dbc *pgxpool.Pool) db.DB {
	return &pg{dbc: dbc}
}

func (p *pg) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("ошибка ScanOneContext: %w", err)
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("ошибка ScanAllContext: %w", err)
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...any) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...any) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...any) pgx.Row {
	logQuery(ctx, q, args...)

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p *pg) Close() {
	p.dbc.Close()
}

func logQuery(_ context.Context, q db.Query, args ...any) {
	log.Println(fmt.Sprintf("sql: %s", q.Name), fmt.Sprintf("query: %s", q.QueryRaw), args, "- sql end.")
}
