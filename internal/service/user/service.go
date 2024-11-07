package user

import (
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository"
	"github.com/NikolosHGW/platform-common/pkg/db"
)

type service struct {
	r         repository.UserRepository
	txManager db.TxManager
}

// NewService - конструктор сервиса юзера.
func NewService(r repository.UserRepository, txManager db.TxManager) *service {
	return &service{
		r:         r,
		txManager: txManager,
	}
}
