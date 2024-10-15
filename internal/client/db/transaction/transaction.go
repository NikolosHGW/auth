package transaction

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/client/db"
	"github.com/NikolosHGW/auth/internal/client/db/pg"
	"github.com/jackc/pgx/v5"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager - конструктор менеджера транзакций.
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	_, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		err = fmt.Errorf("ошибка при старте транзакции")
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("паника восстановлена: %w", err)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = fmt.Errorf("ошибка отката транзакции: %w", errRollback)
			}

			return
		}

		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("ошибка при фиксации транзакции: %w", err)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("не удалось выполнить выполнить код, обёрнутый транзакцией: %w", err)
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

	return m.transaction(ctx, txOpts, f)
}
