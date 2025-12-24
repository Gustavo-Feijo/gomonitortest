package user

type service struct {
	repo *repository
}

func newService(repo *repository) *service {
	return &service{repo}
}

func (s *service) CreateUser(name string) (*User, error) {
	user := &User{Name: name}
	return user, s.repo.Create(user)
}

func (s *service) GetUser(id uint) (*User, error) {
	return s.repo.FindByID(id)
}
