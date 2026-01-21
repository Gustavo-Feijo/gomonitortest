package databaseinfra_test

import (
	"context"
	databaseinfra "gomonitor/internal/infra/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabaseConnection(t *testing.T) {
	t.Parallel()

	dbConn, err := databaseinfra.New(t.Context(), testDbCfg)
	require.NotNil(t, dbConn)
	require.Nil(t, err)

	sqlDB, err := dbConn.DB()
	require.NoError(t, err)

	assert.Equal(t, testDbCfg.MaxOpenConns, sqlDB.Stats().MaxOpenConnections)

	require.NoError(t, sqlDB.Close())
}

func TestNewDatabaseConnectionError(t *testing.T) {
	t.Parallel()

	failDbCfg := *testDbCfg
	failDbCfg.Port = "1"
	dbConn, err := databaseinfra.New(t.Context(), &failDbCfg)
	require.Nil(t, dbConn)
	require.NotNil(t, err)
}

func TestNewDatabaseConnection_ClosesOnPingFailure(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	db, err := databaseinfra.New(ctx, testDbCfg)

	require.Nil(t, db)
	require.Error(t, err)
	assert.ErrorContains(t, err, "failed to ping database")
}
