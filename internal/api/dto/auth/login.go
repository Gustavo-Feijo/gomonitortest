package authdto

import "gomonitor/internal/domain/auth"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (r *LoginRequest) ToDomainInput() auth.LoginInput {
	return auth.LoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"token"`
}

func ToLoginResponse(output *auth.LoginOutput) *LoginResponse {
	return &LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}
}
