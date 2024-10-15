package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/service/user/converter"
	serviceUser "github.com/NikolosHGW/auth/internal/service/user/model"
)

func (s *service) Create(ctx context.Context, user *serviceUser.User) (int64, error) {
	repoUser := converter.UserServiceToUserRepo(user)

	var id int64
	err := s.txManager.ReadCommitted(
		ctx,
		func(ctx context.Context) error {
			txID, err := s.r.Create(ctx, repoUser)
			if err != nil {
				return fmt.Errorf("ошибка при создании с транзакцией: %w", err)
			}

			txUser, err := s.r.GetByID(ctx, txID)
			if err != nil {
				return fmt.Errorf("ошибка при получении с транзакцией: %w", err)
			}
			fmt.Println("получен пользователь в транзакции: ", txUser)
			id = txID

			return nil
		},
	)

	if err != nil {
		return id, fmt.Errorf("ошибка создании: %w", err)
	}
	// TODO: заглушка, убрать потом
	fmt.Println("успешно создан пользователь: ", user.Name, user.Email, user.Password, user.Role)

	return id, nil
}
