package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string // ハッシュ化されたパスワード
	IsAdmin  bool
}
