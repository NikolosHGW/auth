package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/service/user/converter"
	serviceUser "github.com/NikolosHGW/auth/internal/service/user/model"
)

func (s *service) Update(ctx context.Context, serviceUser *serviceUser.User) error {
	err := s.r.Update(ctx, converter.UserServiceToUserRepo(serviceUser))
	if err != nil {
		return fmt.Errorf("ошибка при обновлении: %w", err)
	}

	return nil
}
