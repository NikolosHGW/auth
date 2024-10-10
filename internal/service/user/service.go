package user

import "github.com/NikolosHGW/auth/internal/infrastructure/db/repository"

type service struct {
	r repository.UserRepository
}

// NewService - конструктор сервиса юзера.
func NewService(r repository.UserRepository) *service {
	return &service{
		r: r,
	}
}
