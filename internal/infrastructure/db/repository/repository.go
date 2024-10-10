package repository

import (
	"context"

	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
)

// UserRepository - контракт для репо юзера.
type UserRepository interface {
	Create(context.Context, *userpb.CreateRequest) (int64, error)
	GetByID(context.Context, int64) (*userpb.GetResponse, error)
	UpdateByID(context.Context, *userpb.UpdateRequest) error
	DeleteByID(context.Context, int64) error
}
