package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/api/user/converter"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
)

func (i *Implementation) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.APICreateUserToServiceUser(req))
	if err != nil {
		return nil, fmt.Errorf("ошибка создании: %w", err)
	}
	fmt.Println("успешно: ", req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role)

	return &userpb.CreateResponse{
		Id: id,
	}, nil
}
