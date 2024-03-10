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

// Commit commits the transaction.
func (tx tx) Commit(ctx context.Context) error {
	err := tx.Tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("store: failed to commit transaction: %w", err)
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

		return fmt.Errorf("store: failed to rollback transaction: %w", err)
	}

	return nil
}

// NewReadOnlyTx returns a new transaction that only performs read operations.
func NewReadOnlyTx(ctx context.Context, db *pgxpool.Pool) (pgx.Tx, error) {
	pgxTx, err := db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("store: failed to initialize read only transaction: %w", err)
	}

	return &tx{
		Tx: pgxTx,
	}, nil
}

// NewReadWriteTx returns a new transaction that performs read and write operations.
func NewReadWriteTx(ctx context.Context, db *pgxpool.Pool) (pgx.Tx, error) {
	pgxTx, err := db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, fmt.Errorf("store: failed to initialize read write transaction: %w", err)
	}

	return &tx{
		Tx: pgxTx,
	}, nil
}
