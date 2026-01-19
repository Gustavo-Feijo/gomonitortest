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

	depsSuccess, err := deps.New(t.Context(), cfg, slog.Default())
	assert.Nil(t, err)
	assert.NotNil(t, depsSuccess.DB)

	postgresContainer.Terminate(t.Context())
	depsErr, err := deps.New(t.Context(), cfg, slog.Default())
	assert.NotNil(t, err)
	assert.Nil(t, depsErr)
}
