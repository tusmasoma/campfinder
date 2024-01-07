//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CheckIfUserExists(ctx context.Context, name string, opts ...QueryOptions) (bool, error)
	GetUserByEmail(ctx context.Context, name string, opts ...QueryOptions) (User, error)
	Create(ctx context.Context, user User, opts ...QueryOptions) error
	Update(ctx context.Context, user User, opts ...QueryOptions) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string // ハッシュ化されたパスワード
	IsAdmin  bool
}

// 暗号(Hash)化
func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// 暗号(Hash)と入力された平パスワードの比較
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (ur *userRepository) CheckIfUserExists(ctx context.Context, name string, opts ...QueryOptions) (bool, error) {
	var exists bool
	var executor SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `SELECT EXISTS(SELECT 1 FROM Users WHERE name = ?)`
	err := executor.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string, opts ...QueryOptions) (User, error) {
	var executor SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM User
	WHERE email = ?
	`

	var user User
	err := executor.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
	)
	return user, err
}

func (ur *userRepository) Create(ctx context.Context, user User, opts ...QueryOptions) error {
	var executor SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
    INSERT INTO User (
        id, name, email, password
    )
    VALUES (?, ?, ?, ?)
    `

	_, err := executor.ExecContext(
		ctx,
		query,
		uuid.New(),
		user.Name,
		user.Email,
		user.Password,
	)

	return err
}

func (ur *userRepository) Update(ctx context.Context, user User, opts ...QueryOptions) error {
	var executor SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	UPDATE User SET
	name=?, email=?, password=?
	WHERE id = ?
	`
	_, err := executor.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.ID)
	return err
}
