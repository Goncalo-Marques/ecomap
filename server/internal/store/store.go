package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

// pgIntArray returns a postgres array for the given integer elements.
func pgIntArray(elements []int) []pgtype.Int8 {
	pgElements := make([]pgtype.Int8, len(elements))
	for i, e := range elements {
		pgElements[i] = pgtype.Int8{
			Int64: int64(e),
			Valid: true,
		}
	}

	return pgElements
}

// jsonMarshalGeoJSONGeometryPoint marshals the given GeoJSON into geometry point JSON.
func jsonMarshalGeoJSONGeometryPoint(geoJSON domain.GeoJSON) ([]byte, error) {
	if geoJSON == nil {
		return nil, nil
	}

	var geometry domain.GeoJSONGeometryPoint
	if feature, ok := geoJSON.(domain.GeoJSONFeature); ok {
		if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
			geometry = g
		}
	}

	return json.Marshal(geometry)
}

// listSQLWhere returns an SQL WHERE clause for the specified filter fields, including optional location filters.
func listSQLWhere(fields []string, locationFields []string) string {
	if len(fields) == 0 && len(locationFields) == 0 {
		return ""
	}

	// Construct SQL.
	var sqlFilter string
	if len(fields) != 0 {
		for i, field := range fields {
			fields[i] = field + " ILIKE '%%' || $%d || '%%'"
		}

		sqlFilter = strings.Join(fields, " AND ")
	}

	var sqlLocationFilter string
	if len(locationFields) != 0 {
		for i, field := range locationFields {
			locationFields[i] = field + " ILIKE '%%' || $%d || '%%'"
		}

		sqlLocationFilter = "(" + strings.Join(locationFields, " OR ") + ")"
	}

	sql := " WHERE "
	if len(sqlLocationFilter) == 0 {
		sql += sqlFilter
	} else if len(sqlFilter) == 0 {
		sql += sqlLocationFilter
	} else {
		sql += strings.Join([]string{sqlFilter, sqlLocationFilter}, " AND ")
	}

	// Format parameters.
	sqlParamIndices := make([]any, len(fields)+len(locationFields))
	for i := range fields {
		sqlParamIndices[i] = i + 1
	}
	for i := range locationFields {
		sqlParamIndices[len(fields)+i] = len(fields) + 1
	}

	return fmt.Sprintf(sql, sqlParamIndices...)
}

// listSQLOrder returns an SQL ORDER keyword for the specified field and order.
func listSQLOrder(field string, order domain.PaginationOrder) string {
	o := " ASC"
	if order == domain.PaginationOrderDesc {
		o = " DESC"
	}

	return " ORDER BY " + field + o
}

// listSQLLimitOffset returns an SQL LIMIT and OFFSET clause for the specified limit and offset.
func listSQLLimitOffset(limit domain.PaginationLimit, offset domain.PaginationOffset) string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}
