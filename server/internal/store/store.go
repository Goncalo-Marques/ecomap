package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/goncalo-marques/ecomap/server/internal/config"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/store/tx"
)

// Common failure descriptions.
const (
	descriptionFailedQuery          = "store: failed to query"
	descriptionFailedExec           = "store: failed to exec"
	descriptionFailedScanRow        = "store: failed to scan row"
	descriptionFailedScanRows       = "store: failed to scan rows"
	descriptionFailedMarshalGeoJSON = "store: failed to marshal geojson"
)

// migrationsURL defines the source url of the migrations.
const migrationsURL = "file://database/migrations"

// store defines the store structure.
type store struct {
	database *pgxpool.Pool
}

// New returns a new store.
func New(ctx context.Context, config config.Database) (*store, error) {
	// Initialize database connection pool.
	database, err := pgxpool.New(ctx, config.URL)
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
		database: database,
	}, nil
}

// NewTx returns a new transaction for the current database.
func (s *store) NewTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode) (pgx.Tx, error) {
	return tx.New(ctx, s.database, isoLevel, accessMode)
}

// Close closes the database.
func (s *store) Close() {
	s.database.Close()
}

// constraintNameFromError returns the name of the constraint of the given error. If the error is not of type
// pgconn.PgError, an empty string is returned.
func constraintNameFromError(err error) string {
	if pqErr, ok := err.(*pgconn.PgError); ok {
		return pqErr.ConstraintName
	}

	return ""
}

// listSQLWhere returns a SQL WHERE clause for the specified filter fields using the specified logical operator.
func listSQLWhere(fields []string, logicalOperator domain.PaginationLogicalOperator) string {
	if len(fields) == 0 {
		return ""
	}

	// Construct SQL.
	for i, field := range fields {
		fields[i] = field + " ILIKE '%%' || $%d || '%%'"
	}

	lo := " AND "
	if logicalOperator == domain.PaginationLogicalOperatorOr {
		lo = " OR "
	}

	sql := " WHERE " + strings.Join(fields, lo)

	// Format parameters.
	sqlParamIndices := make([]any, len(fields))
	for i := range fields {
		sqlParamIndices[i] = i + 1
	}

	return fmt.Sprintf(sql, sqlParamIndices...)
}

// listSQLOrder returns a SQL ORDER keyword for the specified field and order.
func listSQLOrder(field string, order domain.PaginationOrder) string {
	o := " ASC"
	if order == domain.PaginationOrderDesc {
		o = " DESC"
	}

	return " ORDER BY " + field + o
}

// listSQLLimitOffset returns a SQL LIMIT and OFFSET clause for the specified limit and offset.
func listSQLLimitOffset(limit domain.PaginationLimit, offset domain.PaginationOffset) string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}
