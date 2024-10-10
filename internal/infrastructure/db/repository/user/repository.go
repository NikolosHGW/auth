package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	pgxCon *pgxpool.Pool
}

// NewUser - конструктор репозитория юзера.
func NewUser(pgxCon *pgxpool.Pool) *userRepo {
	return &userRepo{pgxCon: pgxCon}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	query := `INSERT INTO users (name, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.
		pgxCon.
		QueryRow(
			ctx,
			query,
			user.Name,
			user.Password,
			user.Email,
			user.Role,
		).
		Scan(&id)

	if err != nil {
		return id, fmt.Errorf("ошибка при сохранении пользователя: %w", err)
	}

	return id, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, password, email, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.pgxCon.QueryRow(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3
	`
	_, err := r.pgxCon.Exec(ctx, query, user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	return nil
}

func (r *userRepo) DeleteByID(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.pgxCon.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	return nil
}
