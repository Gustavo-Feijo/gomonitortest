package authdto_test

import (
	authdto "gomonitor/internal/api/dto/auth"
	"gomonitor/internal/domain/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDto_LoginRequest(t *testing.T) {
	loginRequest := &authdto.LoginRequest{
		Email:    "test@test.com",
		Password: "test123",
	}

	expectedLoginInput := auth.LoginInput{
		Email:    "test@test.com",
		Password: "test123",
	}

	loginInput := loginRequest.ToDomainInput()

	assert.EqualValues(t, expectedLoginInput, loginInput)
}

func TestDto_LoginResponse(t *testing.T) {
	loginOutput := &auth.LoginOutput{
		RefreshToken: "testRefreshToken",
		AccessToken:  "testAccessToken",
	}

	expectedLoginResponse := &authdto.LoginResponse{
		RefreshToken: "testRefreshToken",
		AccessToken:  "testAccessToken",
	}

	loginResponse := authdto.ToLoginResponse(loginOutput)

	assert.EqualValues(t, expectedLoginResponse, loginResponse)
}
