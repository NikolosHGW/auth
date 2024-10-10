package service

import (
	"context"

	serviceUser "github.com/NikolosHGW/auth/internal/service/user/model"
)

// UserService - контракт для сервиса юзера.
type UserService interface {
	Create(ctx context.Context, user *serviceUser.User) (int64, error)
	Get(ctx context.Context, id int64) (*serviceUser.User, error)
	Update(ctx context.Context, serviceUser *serviceUser.User) error
	Delete(ctx context.Context, id int64) error
}
