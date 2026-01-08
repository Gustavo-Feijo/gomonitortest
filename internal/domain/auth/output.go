package auth

type LoginOutput struct {
	RefreshToken string
	AccessToken  string
}

type RefreshOutput struct {
	AccessToken string
}
