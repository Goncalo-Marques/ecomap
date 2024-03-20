//go:build integration

package database_test

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"

	"github.com/goncalo-marques/ecomap/server/test/container"
)

// migrationsURL defines the source url of the migrations.
const migrationsURL = "file://migrations"

func TestMigrations(t *testing.T) {
	ctx := context.Background()

	databaseContainer := container.NewDatabase(ctx)
	defer databaseContainer.Terminate(ctx)

	m, err := migrate.New(migrationsURL, databaseContainer.ConnectionString(ctx))
	require.NoError(t, err)
	defer m.Close()

	err = m.Up()
	require.NoError(t, err)

	err = m.Down()
	require.NoError(t, err)
}
