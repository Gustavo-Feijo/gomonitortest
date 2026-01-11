package authdto_test

import (
	authdto "gomonitor/internal/api/dto/auth"
	"gomonitor/internal/domain/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDto_RefreshRequest(t *testing.T) {
	refreshRequest := &authdto.RefreshRequest{
		"testRefreshToken",
	}

	expectedRefreshInput := auth.RefreshInput{
		RefreshToken: "testRefreshToken",
	}

	refreshInput := refreshRequest.ToDomainInput()

	assert.EqualValues(t, expectedRefreshInput, refreshInput)
}

func TestDto_RefreshResponse(t *testing.T) {
	refreshOutput := &auth.RefreshOutput{
		AccessToken: "testAccessToken",
	}

	expectedRefreshResponse := &authdto.RefreshResponse{
		AccessToken: "testAccessToken",
	}

	refreshResponse := authdto.ToRefreshResponse(refreshOutput)

	assert.EqualValues(t, expectedRefreshResponse, refreshResponse)
}
