package user

import "context"

type service struct {
	repo *repository
}

func newService(repo *repository) *service {
	return &service{repo}
}

func (s *service) CreateUser(ctx context.Context, name string) (*User, error) {
	user := &User{Name: name}
	return user, s.repo.Create(ctx, user)
}

func (s *service) GetUser(ctx context.Context, id uint) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
