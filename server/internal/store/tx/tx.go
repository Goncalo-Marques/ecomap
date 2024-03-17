package tx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// tx defines the transaction structure.
type tx struct {
	pgx.Tx
}

// New returns a new transaction.
func New(ctx context.Context, database *pgxpool.Pool, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode) (pgx.Tx, error) {
	pgxTx, err := database.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   isoLevel,
		AccessMode: accessMode,
	})
	if err != nil {
		return nil, fmt.Errorf("tx: failed to initialize transaction: %w", err)
	}

	return &tx{
		Tx: pgxTx,
	}, nil
}

// Commit commits the transaction.
func (tx tx) Commit(ctx context.Context) error {
	err := tx.Tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx: failed to commit transaction: %w", err)
	}

	return nil
}

// Rollback rolls back the transaction.
// It is safe to call this method after the Commit.
func (tx tx) Rollback(ctx context.Context) error {
	err := tx.Tx.Rollback(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrTxClosed) {
			return nil
		}

		return fmt.Errorf("tx: failed to rollback transaction: %w", err)
	}

	return nil
}
