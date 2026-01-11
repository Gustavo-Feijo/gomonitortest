package userdto_test

import (
	userdto "gomonitor/internal/api/dto/user"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDto_CreateUserRequest(t *testing.T) {
	role := testutil.Ptr(identity.RoleUser)
	createUserRequest := &userdto.CreateUserRequest{
		Name:     "test",
		Email:    "test@test.com",
		Password: "test123",
		UserName: "test",
		Role:     role,
	}

	expectedCreateUserInput := user.CreateUserInput{
		Name:     "test",
		Email:    "test@test.com",
		Password: "test123",
		UserName: "test",
		Role:     role,
	}

	createUserInput := createUserRequest.ToDomainInput()

	assert.EqualValues(t, expectedCreateUserInput, createUserInput)
}

func TestDto_CreateUserResponse(t *testing.T) {
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

	expectedCreateUserResponse := &userdto.CreateUserResponse{
		ID:        1,
		Name:      "test",
		Email:     "test@test.com",
		UserName:  "test",
		Role:      identity.RoleUser,
		CreatedAt: now,
	}

	createUserResponse := userdto.ToCreateUserResponse(user)

	assert.EqualValues(t, expectedCreateUserResponse, createUserResponse)
}
