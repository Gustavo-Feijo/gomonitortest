package password

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}

type passwordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) PasswordHasher {
	return &passwordHasher{cost: cost}
}

func (ps *passwordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), ps.cost)
	return string(bytes), err
}

func (ps *passwordHasher) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
