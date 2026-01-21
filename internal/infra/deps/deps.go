package deps

import (
	"context"
	"errors"
	"fmt"
	"gomonitor/internal/config"
	databaseinfra "gomonitor/internal/infra/database"
	redisinfra "gomonitor/internal/infra/redis"
	"gomonitor/internal/pkg/jwt"
	"gomonitor/internal/pkg/password"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Dependencies for the service.
type Deps struct {
	DB           *gorm.DB
	Hasher       password.PasswordHasher
	Logger       *slog.Logger
	Redis        redisinfra.RedisClient
	TokenManager jwt.TokenManager
}

// New creates the necessary instances.
func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Deps, func(ctx context.Context) error, error) {
	db, err := databaseinfra.New(ctx, cfg.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("error at opening db conn: %w", err)
	}

	// Treat redis connection. Redis is optional for full functionality.
	rdb := redisinfra.New(ctx, cfg.Redis, cfg.CircuitBreaker, logger)

	cleanup := func(ctx context.Context) error {
		var errs []error

		if rdb != nil {
			if err := rdb.Close(); err != nil {
				errs = append(errs, fmt.Errorf("redis close: %w", err))
			}
		}

		if db != nil {
			sqlDb, _ := db.DB()
			if err := sqlDb.Close(); err != nil {
				errs = append(errs, fmt.Errorf("db close: %w", err))
			}
		}

		return errors.Join(errs...)
	}

	return &Deps{
		DB:           db,
		Hasher:       password.NewPasswordHasher(bcrypt.DefaultCost),
		Logger:       logger,
		Redis:        rdb,
		TokenManager: jwt.NewTokenManager(cfg.Auth),
	}, cleanup, nil
}
