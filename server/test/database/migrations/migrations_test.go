//go:build integration

package migrations_test

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"

	"github.com/goncalo-marques/ecomap/server/test/container"
)

// migrationsURL defines the source url of the migrations.
const migrationsURL = "file://database/migrations"

func TestMigrations(t *testing.T) {
	ctx := context.Background()

	databaseContainer := container.NewDatabase(ctx)
	defer databaseContainer.Terminate(ctx)

	m, err := migrate.New(migrationsURL, databaseContainer.ConnectionString())
	require.NoError(t, err)
	defer m.Close()

	err = m.Up()
	require.NoError(t, err)

	err = m.Down()
	require.NoError(t, err)
}
