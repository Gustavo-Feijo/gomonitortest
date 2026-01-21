package app

import (
	"gomonitor/internal/config"
	"gomonitor/internal/testutil"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHandler struct {
	called bool
}

func (f *fakeHandler) RegisterRoutes(r *gin.RouterGroup) {
	f.called = true
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})
}

func TestNewApp(t *testing.T) {
	testutil.StartTestDB(t)
	testutil.StartTestRedis(t)

	cfg, err := config.Load()
	require.Nil(t, err)
	require.NotNil(t, cfg)

	app, cleanup, err := New(t.Context(), cfg, slog.Default())
	require.Nil(t, err)
	require.NotNil(t, app)

	assert.Nil(t, cleanup(t.Context()))
}

func TestNewAppDepsError(t *testing.T) {
	postgresContainer := testutil.StartTestDB(t)
	if err := postgresContainer.Terminate(t.Context()); err != nil {
		t.Logf("cleanup: failed to terminate postgres container: %v", err)
	}
	testutil.StartTestRedis(t)

	cfg, err := config.Load()
	require.Nil(t, err)
	require.NotNil(t, cfg)

	app, cleanup, err := New(t.Context(), cfg, slog.Default())
	assert.NotNil(t, err)
	assert.Nil(t, app)

	assert.Nil(t, cleanup)
}

func TestNewAppMigrationError(t *testing.T) {
	testutil.StartTestDB(t)
	testutil.StartTestRedis(t)

	cfg, err := config.Load()
	require.Nil(t, err)
	require.NotNil(t, cfg)

	cfg.Database.MigrationsPath = "nonexistent"

	app, cleanup, err := New(t.Context(), cfg, slog.Default())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "error running migrations")
	assert.Nil(t, app)

	assert.Nil(t, cleanup)
}

func TestRegisterRoutes_CallsHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := &fakeHandler{}

	registerRoutes(r, h)

	assert.True(t, h.called)

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRegisterRoutes_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	registerRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}
