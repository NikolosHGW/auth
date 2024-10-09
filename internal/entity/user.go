package entity

import "time"

// User - модель пользователя.
type User struct {
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	ID        int64      `db:"id"`
	Role      int        `db:"role"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
