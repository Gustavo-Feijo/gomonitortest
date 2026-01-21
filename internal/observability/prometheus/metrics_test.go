package prometheus_test

import (
	pkgprometheus "gomonitor/internal/observability/prometheus"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestInitRegistersMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()

	require.NotPanics(t, func() {
		pkgprometheus.Init(reg)
	})
}
