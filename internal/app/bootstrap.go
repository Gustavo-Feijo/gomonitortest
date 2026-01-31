package app

import (
	"context"
	"gomonitor/internal/container"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/observability/logging"
	"gomonitor/internal/pkg/identity"
	"log/slog"

	"gorm.io/gorm"
)

func bootstrapApp(ctx context.Context, container *container.Container) error {
	return container.Deps.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := createAdminUser(ctx, tx, container); err != nil {
			return err
		}

		return nil
	})
}

// createAdminUser initializes the first user on the application as admin.
func createAdminUser(ctx context.Context, tx *gorm.DB, c *container.Container) error {
	logger := c.Deps.Logger
	hasher := c.Deps.Hasher
	cfg := c.Cfg.Admin

	// Create new repository and service instead of using the container one due to the transactional nature.
	userRepo := c.Repositories.User.WithTx(tx)

	count, err := userRepo.Count(ctx)
	if err != nil {
		logger.Error("error getting user count", slog.Any("err", err))
		return err
	}

	if count != 0 {
		return nil
	}

	userSvcDeps := &user.ServiceDeps{
		Hasher:   hasher,
		Logger:   logger,
		UserRepo: userRepo,
	}
	userSvc := user.NewService(userSvcDeps)

	principal := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleAdmin,
		Source: identity.AuthInternal,
	}

	ctxWithLogging := logging.WithContext(ctx, logger)
	internalCtx := identity.WithPrincipal(ctxWithLogging, principal)

	var role = identity.RoleAdmin
	adminUser := user.CreateUserInput{
		Name:     "admin",
		UserName: "admin",
		Email:    cfg.Email,
		Password: cfg.Password,
		Role:     &role,
	}

	newUser, err := userSvc.CreateUser(internalCtx, adminUser)
	if err != nil {
		logger.Error("error creating admin user", slog.Any("err", err))
		return err
	}

	logger.Info(
		"success creating admin user",
		slog.Uint64("user_id", uint64(newUser.ID)),
		slog.String("user_email", newUser.Email),
	)

	return nil
}
