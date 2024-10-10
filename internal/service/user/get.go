package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/service/user/converter"
	serviceUser "github.com/NikolosHGW/auth/internal/service/user/model"
)

func (s *service) Get(ctx context.Context, id int64) (*serviceUser.User, error) {
	repoUser, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return converter.UserRepoToServiceRepo(repoUser), nil
}
