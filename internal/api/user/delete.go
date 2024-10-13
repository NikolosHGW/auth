package user

import (
	"context"
	"fmt"

	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete удаляет пользователя.
func (i *Implementation) Delete(ctx context.Context, req *userpb.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при удалении: %w", err)
	}
	fmt.Println(req.Id)

	return &emptypb.Empty{}, nil
}
