package redisinfra

import (
	"context"
	"gomonitor/internal/config"
	"gomonitor/internal/testutil"
	"log"
	"os"
	"testing"
)

var (
	testRedisCfg *config.RedisConfig
	testCbCfg    *config.CircuitBreakerConfig
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	_, addr, cleanup, err := testutil.StartRedis(ctx)
	if err != nil {
		log.Fatal(err)
	}

	testRedisCfg = &config.RedisConfig{
		Addr:     addr,
		Database: 0,
	}

	testCbCfg = &config.CircuitBreakerConfig{
		MaxRequests: 5,
		MaxFailures: 5,
	}

	code := m.Run()
	_ = cleanup(ctx)
	os.Exit(code)
}
