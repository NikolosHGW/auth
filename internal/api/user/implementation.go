package user

import (
	"github.com/NikolosHGW/auth/internal/service"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
)

type implementation struct {
	userpb.UserV1Server
	userService service.UserService
}

// NewImplementation - конструктор gRPC сервера для user.
func NewImplementation(userService service.UserService) *implementation {
	return &implementation{userService: userService}
}
