package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/api/user/converter"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
)

// Get ищет пользователя в бд.
func (i *Implementation) Get(ctx context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
	serviceUser, err := i.userService.Get(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения: %w", err)
	}
	fmt.Println("успешно: ", req.Id)

	return converter.ServiceUserToAPIGetUser(serviceUser), nil
}
