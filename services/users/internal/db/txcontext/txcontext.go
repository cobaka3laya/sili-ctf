package txcontext

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Manager struct {
	pool *pgxpool.Pool
}

func (m *Manager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)

	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err = fn(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		tx.Rollback(ctx)
	}

	return err
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
