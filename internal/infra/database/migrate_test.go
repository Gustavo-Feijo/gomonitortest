package databaseinfra_test

import (
	"context"
	databaseinfra "gomonitor/internal/infra/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMigrations(t *testing.T) {
	dbConn, err := databaseinfra.New(t.Context(), testDbCfg)
	require.Nil(t, err)
	require.NotNil(t, dbConn)

	err = databaseinfra.RunMigrations(t.Context(), testDbCfg, dbConn)
	require.Nil(t, err)
}

func TestRunMigrationsCtxErr(t *testing.T) {
	dbConn, err := databaseinfra.New(t.Context(), testDbCfg)
	require.Nil(t, err)
	require.NotNil(t, dbConn)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err = databaseinfra.RunMigrations(ctx, testDbCfg, dbConn)
	require.NotNil(t, err)
}

func TestRunMigrationsClosedConn(t *testing.T) {
	dbConn, err := databaseinfra.New(t.Context(), testDbCfg)
	require.Nil(t, err)
	require.NotNil(t, dbConn)

	sqlDb, _ := dbConn.DB()
	err = sqlDb.Close()
	require.Nil(t, err)

	err = databaseinfra.RunMigrations(t.Context(), testDbCfg, dbConn)
	require.NotNil(t, err)
	assert.ErrorContains(t, err, "could not create migration driver")
}

func TestRunMigrationsWrongPath(t *testing.T) {
	wrongPathDbConfig := *testDbCfg
	wrongPathDbConfig.MigrationsPath = "wrongpath"
	dbConn, err := databaseinfra.New(t.Context(), &wrongPathDbConfig)
	require.Nil(t, err)
	require.NotNil(t, dbConn)

	err = databaseinfra.RunMigrations(t.Context(), &wrongPathDbConfig, dbConn)
	require.NotNil(t, err)
	assert.ErrorContains(t, err, "could not create migrate instance")
}
