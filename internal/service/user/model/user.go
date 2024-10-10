package model

import (
	"database/sql"
	"time"
)

// User - модель пользователя для сервиса.
type User struct {
	Name      string
	Email     string
	Password  string
	ID        int64
	Role      int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
