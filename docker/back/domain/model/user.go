package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Email    string    `db:"email"`
	Password string    `db:"password"` // ハッシュ化されたパスワード
	IsAdmin  bool      `db:"is_admin"`
}
