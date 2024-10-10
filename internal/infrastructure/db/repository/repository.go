package repository

import (
	"context"

	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user/model"
)

// UserRepository - контракт для репо юзера.
type UserRepository interface {
	Create(context.Context, *model.User) (int64, error)
	GetByID(context.Context, int64) (*model.User, error)
	Update(context.Context, *model.User) error
	DeleteByID(context.Context, int64) error
}
