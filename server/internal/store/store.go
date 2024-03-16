package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/goncalo-marques/ecomap/server/internal/config"
	"github.com/goncalo-marques/ecomap/server/internal/store/tx"
)

// migrationsURL defines the source url of the migrations.
const migrationsURL = "file://db/migrations"

// store defines the store structure.
type store struct {
	db *pgxpool.Pool
}

// New returns a new store.
func New(ctx context.Context, config config.Database) (*store, error) {
	// Initialize database connection pool.
	db, err := pgxpool.New(ctx, config.URL)
	if err != nil {
		return nil, fmt.Errorf("store: failed to initialize pool: %w", err)
	}

	// Apply migrations based on the configuration.
	if config.Migrations.Apply {
		m, err := migrate.New(migrationsURL, config.URL)
		if err != nil {
			return nil, fmt.Errorf("store: failed to initialize migrate: %w", err)
		}
		defer m.Close()

		if err := m.Migrate(config.Migrations.Version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("store: failed to apply migrations: %w", err)
		}
	}

	return &store{
		db: db,
	}, nil
}

// NewTx returns a new transaction for the current db.
func (s *store) NewTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode) (pgx.Tx, error) {
	return tx.New(ctx, s.db, isoLevel, accessMode)
}

// Close closes the db.
func (s *store) Close() {
	s.db.Close()
}
