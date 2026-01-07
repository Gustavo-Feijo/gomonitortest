package auth

type LoginOutput struct {
	RefreshToken string
	ApiToken     string
}

type RefreshOutput struct {
	ApiToken string
}
