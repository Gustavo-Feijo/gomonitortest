package deps

import (
	"fmt"
	"gomonitor/internal/config"
	databaseinfra "gomonitor/internal/infra/database"
	redisinfra "gomonitor/internal/infra/redis"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Dependencies for the service.
type Deps struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewDeps creates the necessary instances.
func NewDeps(cfg *config.Config) (*Deps, error) {
	db, err := databaseinfra.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("error at opening db conn: %w", err)
	}

	return &Deps{
		DB:    db,
		Redis: redisinfra.New(cfg.Redis),
	}, nil
}
