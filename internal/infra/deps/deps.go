package deps

import (
	"context"
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
	Redis        *redisinfra.RedisClient
	TokenManager jwt.TokenManager
}

// NewDeps creates the necessary instances.
func NewDeps(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Deps, error) {
	db, err := databaseinfra.New(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("error at opening db conn: %w", err)
	}

	// Treat redis connection. Redis is optional for full functionality.
	rdb, err := redisinfra.New(ctx, cfg.Redis, logger)
	if err != nil {
		logger.Error("failed to initialize redis connection", slog.Any("err", err))
	}

	return &Deps{
		DB:           db,
		Hasher:       password.NewPasswordHasher(bcrypt.DefaultCost),
		Logger:       logger,
		Redis:        rdb,
		TokenManager: jwt.NewTokenManager(cfg.Auth),
	}, nil
}
