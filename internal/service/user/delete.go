package user

import (
	"context"
	"fmt"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.r.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении: %w", err)
	}

	return nil
}
