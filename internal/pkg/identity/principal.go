package identity

import (
	"context"

	"github.com/google/uuid"
)

// External or internal authentication.
type AuthSource string

const (
	AuthExternal AuthSource = "external"
	AuthInternal AuthSource = "internal"
)

type Principal struct {
	UserID uint
	Role   UserRole
	Source AuthSource

	JTI        *uuid.UUID // nil for access tokens
	RefreshJTI *uuid.UUID
}

type principalKeyType struct{}

var principalKey = principalKeyType{}

func WithPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}

func PrincipalFromContext(ctx context.Context) (*Principal, bool) {
	p, ok := ctx.Value(principalKey).(*Principal)
	return p, ok
}
