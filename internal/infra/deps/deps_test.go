package deps_test

import (
	"gomonitor/internal/config"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/testutil"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeps(t *testing.T) {
	postgresContainer := testutil.StartTestDB(t)
	testutil.StartTestRedis(t)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("couldn't load config: %v", err)
	}

	depsSuccess, _, err := deps.New(t.Context(), cfg, slog.Default())
	assert.Nil(t, err)
	assert.NotNil(t, depsSuccess.DB)

	if err := postgresContainer.Terminate(t.Context()); err != nil {
		t.Logf("cleanup: failed to terminate postgres container: %v", err)
	}

	depsErr, _, err := deps.New(t.Context(), cfg, slog.Default())
	assert.NotNil(t, err)
	assert.Nil(t, depsErr)
}

func TestNewDepsCleanup(t *testing.T) {
	testutil.StartTestDB(t)
	testutil.StartTestRedis(t)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("couldn't load config: %v", err)
	}

	_, cleanup, err := deps.New(t.Context(), cfg, slog.Default())
	assert.Nil(t, err)
	assert.NotNil(t, cleanup)

	err = cleanup(t.Context())
	assert.Nil(t, err)
}
