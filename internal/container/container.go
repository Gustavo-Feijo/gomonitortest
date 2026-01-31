package container

import (
	authhandler "gomonitor/internal/api/handlers/auth"
	userhandler "gomonitor/internal/api/handlers/user"
	"gomonitor/internal/config"
	"gomonitor/internal/domain/auth"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/pkg/ratelimit"
)

type Container struct {
	Deps *deps.Deps
	Cfg  *config.Config

	RateLimiters RateLimiters
	Repositories *Repositories
	Services     *Services
	Handler      *Handlers
}

type RateLimiters struct {
	IPLimiter ratelimit.RateLimiter
}

type Repositories struct {
	User         user.UserRepository
	RefreshToken auth.RefreshTokenRepository
}

type Services struct {
	Auth auth.Service
	User user.Service
}

type Handlers struct {
	Auth *authhandler.Handler
	User *userhandler.Handler
}

func New(deps *deps.Deps, cfg *config.Config) *Container {
	c := &Container{
		Deps:         deps,
		Cfg:          cfg,
		Repositories: &Repositories{},
		Services:     &Services{},
		Handler:      &Handlers{},
	}

	c.RateLimiters.IPLimiter = ratelimit.New(
		ratelimit.WithLimiter(
			ratelimit.NewRedisLimiter(
				deps.Redis,
				ratelimit.WithLimit(cfg.RateLimit.IPLimit),
				ratelimit.WithPrefix("rate_limit"),
				ratelimit.WithWindow(cfg.RateLimit.IPWindow),
			),
		),
		ratelimit.WithFallback(
			ratelimit.NewMemoryLimiter(
				ratelimit.WithLimit(cfg.RateLimit.IPLimit),
				ratelimit.WithPrefix("rate_limit"),
				ratelimit.WithWindow(cfg.RateLimit.IPWindow),
			),
		),
	)

	c.Repositories.User = user.NewUserRepository(deps.DB)
	c.Repositories.RefreshToken = auth.NewRefreshTokenRepository(deps.DB)

	c.Services.Auth = auth.NewService(&auth.ServiceDeps{
		AuthConfig:       cfg.Auth,
		Hasher:           deps.Hasher,
		Logger:           deps.Logger,
		RefreshTokenRepo: c.Repositories.RefreshToken,
		UserRepo:         c.Repositories.User,
		TokenManager:     deps.TokenManager,
	})

	c.Services.User = user.NewService(&user.ServiceDeps{
		Hasher:   deps.Hasher,
		UserRepo: c.Repositories.User,
		Logger:   deps.Logger,
	})

	c.Handler.Auth = authhandler.NewHandler(deps.Logger, c.Services.Auth, deps.TokenManager)
	c.Handler.User = userhandler.NewHandler(deps.Logger, c.Services.User, deps.TokenManager)

	return c
}
