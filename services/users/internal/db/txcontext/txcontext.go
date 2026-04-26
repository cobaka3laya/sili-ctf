package txcontext

import (
	"context"
	"fmt"

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
		return fmt.Errorf("txmanager: starting transaction failed: %w", err)
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err = fn(ctx); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("txmanager: error during transaction: %w", err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("txmanager: commit failed: %w", err)
	}

	return nil
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
