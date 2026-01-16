package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPrincipal(t *testing.T) {
	principal := &Principal{
		UserID: 1,
		Role:   RoleAdmin,
		Source: AuthExternal,
	}
	ctx := WithPrincipal(t.Context(), principal)
	assert.NotNil(t, ctx)

	p, ok := ctx.Value(principalKey).(*Principal)

	assert.True(t, ok)
	assert.NotNil(t, p)
}

func TestPrincipalFromContext(t *testing.T) {
	principal := &Principal{
		UserID: 1,
		Role:   RoleAdmin,
		Source: AuthExternal,
	}
	ctx := WithPrincipal(t.Context(), principal)

	p, ok := PrincipalFromContext(ctx)

	assert.True(t, ok)
	assert.NotNil(t, p)
}
