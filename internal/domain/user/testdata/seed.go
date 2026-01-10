package testdata

import (
	"fmt"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/pkg/identity"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// password123
const TestPasswordHash = "$2a$08$XK3WvA63SGoqzBvYLDBH8uZttWnDNzR8AYXMBtk.ZgtE6TdodbAyS"

func SeedUser(t *testing.T, db *gorm.DB, index int) *user.User {
	t.Helper()

	u := &user.User{
		Name:     fmt.Sprintf("Test User %d", index),
		UserName: fmt.Sprintf("test%d", index),
		Email:    fmt.Sprintf("test%d@test.com", index),
		Password: TestPasswordHash,
		Role:     identity.RoleUser,
	}

	err := db.WithContext(t.Context()).Create(u).Error
	require.NoError(t, err)

	return u
}

func SeedUsers(t *testing.T, db *gorm.DB, n int) []*user.User {
	t.Helper()

	users := make([]*user.User, 0, n)

	for i := 0; i < n; i++ {
		u := SeedUser(t, db, i)
		users = append(users, u)
	}

	return users
}
