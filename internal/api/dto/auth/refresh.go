package authdto

import "gomonitor/internal/domain/auth"

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *RefreshRequest) ToDomainInput() auth.RefreshInput {
	return auth.RefreshInput{
		RefreshToken: r.RefreshToken,
	}
}

type RefreshResponse struct {
	AccessToken string `json:"token"`
}

func ToRefreshResponse(output *auth.RefreshOutput) *RefreshResponse {
	return &RefreshResponse{
		AccessToken: output.AccessToken,
	}
}
