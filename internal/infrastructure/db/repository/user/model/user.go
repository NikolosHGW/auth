package model

import (
	"database/sql"
	"time"
)

// User - модель пользователя.
type User struct {
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	ID        int64        `db:"id"`
	Role      int32        `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
