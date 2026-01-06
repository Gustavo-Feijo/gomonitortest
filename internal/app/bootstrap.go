package app

import (
	"context"
	"gomonitor/internal/config"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/pkg/identity"
	"log/slog"

	"gorm.io/gorm"
)

func bootstrapApp(ctx context.Context, cfg *config.Config, deps *deps.Deps) error {
	return deps.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := createAdminUser(ctx, cfg.Admin, tx, deps.Logger); err != nil {
			return err
		}

		return nil
	})
}

// createAdminUser initializes the first user on the application as admin.
func createAdminUser(ctx context.Context, cfg *config.AdminConfig, db *gorm.DB, logger *slog.Logger) error {
	userRepo := user.NewRepository(db)

	count, err := userRepo.Count(ctx)
	if err != nil {
		logger.Error("error getting user count", slog.Any("err", err))
		return err
	}

	if count != 0 {
		return nil
	}

	userSvcDeps := &user.ServiceDeps{
		Logger: logger,
		DB:     db,
	}
	userSvc := user.NewService(userSvcDeps)

	principal := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleAdmin,
		Source: identity.AuthInternal,
	}

	internalCtx := identity.WithPrincipal(ctx, principal)

	var role identity.UserRole = identity.RoleAdmin
	adminUser := user.CreateUserRequest{
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
