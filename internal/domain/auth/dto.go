package auth

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginResponse struct {
	RefreshToken string `json:"refresh_token"`
	ApiToken     string `json:"token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	ApiToken string `json:"token"`
}
