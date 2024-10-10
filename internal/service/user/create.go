package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/service/user/converter"
	serviceUser "github.com/NikolosHGW/auth/internal/service/user/model"
)

func (s *service) Create(ctx context.Context, user *serviceUser.User) (int64, error) {
	repoUser := converter.UserServiceToUserRepo(user)
	id, err := s.r.Create(ctx, repoUser)
	if err != nil {
		return id, fmt.Errorf("ошибка создании: %w", err)
	}
	// TODO: заглушка, убрать потом
	fmt.Println("успешно создан пользователь: ", user.Name, user.Email, user.Password, user.Role)

	return id, nil
}
