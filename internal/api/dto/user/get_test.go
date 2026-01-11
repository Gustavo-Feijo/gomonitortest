package userdto_test

import (
	userdto "gomonitor/internal/api/dto/user"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/pkg/identity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDto_GetUserRequest(t *testing.T) {
	getUserRequest := &userdto.GetUserRequest{
		ID: uint(1),
	}

	expectedGetUserInput := user.GetUserInput{
		ID: uint(1),
	}

	getUserInput := getUserRequest.ToDomainInput()

	assert.EqualValues(t, expectedGetUserInput, getUserInput)
}

func TestDto_GetUserResponse(t *testing.T) {
	now := time.Now()

	user := &user.User{
		ID:        1,
		Name:      "test",
		Email:     "test@test.com",
		UserName:  "test",
		Password:  "test",
		Role:      identity.RoleUser,
		CreatedAt: now,
		UpdatedAt: now,
	}

	expectedGetUserResponse := &userdto.GetUserResponse{
		ID:        1,
		Name:      "test",
		Email:     "test@test.com",
		UserName:  "test",
		Role:      identity.RoleUser,
		CreatedAt: now,
		UpdatedAt: now,
	}

	getUserResponse := userdto.ToGetUserResponse(user)

	assert.EqualValues(t, expectedGetUserResponse, getUserResponse)
}
