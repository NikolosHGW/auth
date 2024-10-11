package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/api/user/converter"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *implementation) Update(ctx context.Context, req *userpb.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, converter.APIUpdateUserToServiceUser(req))
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении: %w", err)
	}
	fmt.Println(req.Id, req.Name, req.Email)

	return &emptypb.Empty{}, nil
}
