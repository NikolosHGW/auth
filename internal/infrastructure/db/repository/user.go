package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NikolosHGW/auth/internal/entity"
	"github.com/jackc/pgx/v5"
)

type userRepo struct {
	pgxCon *pgx.Conn
}

// NewUser - конструктор репозитория юзера.
func NewUser(pgxCon *pgx.Conn) *userRepo {
	return &userRepo{pgxCon: pgxCon}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
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
		Scan(&user.ID)

	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении пользователя: %w", err)
	}

	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, name, password, email, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.pgxCon.QueryRow(ctx, query, id)

	var updatedAt sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	} else {
		user.UpdatedAt = nil
	}

	return &user, nil
}

func (r *userRepo) UpdateByID(ctx context.Context, user *entity.User) error {
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
