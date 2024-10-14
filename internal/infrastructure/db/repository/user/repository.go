package user

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/auth/internal/client/db"
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user/model"
)

const repositoryName = "user_repository"

type userRepo struct {
	db db.Client
}

// NewUser - конструктор репозитория юзера.
func NewUser(db db.Client) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	q := db.Query{
		Name:     repositoryName + ".Create",
		QueryRaw: `INSERT INTO users (name, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id`,
	}
	err := r.
		db.DB().
		QueryRowContext(
			ctx,
			q,
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
	q := db.Query{
		Name:     repositoryName + ".GetByID",
		QueryRaw: `SELECT id, name, password, email, role, created_at, updated_at FROM users WHERE id = $1`,
	}
	err := r.db.DB().ScanOneContext(ctx, &user, q, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	q := db.Query{
		Name: repositoryName + ".Update",
		QueryRaw: `
			UPDATE users
			SET name = $1, email = $2
			WHERE id = $3
		`,
	}

	_, err := r.db.DB().ExecContext(ctx, q, user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	return nil
}

func (r *userRepo) DeleteByID(ctx context.Context, id int64) error {
	q := db.Query{
		Name:     repositoryName + ".Delete",
		QueryRaw: "DELETE FROM users WHERE id = $1",
	}
	_, err := r.db.DB().ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	return nil
}
