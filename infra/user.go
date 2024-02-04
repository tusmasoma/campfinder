//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package infra

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CheckIfUserExists(
	ctx context.Context,
	email string,
	opts ...repository.QueryOptions,
) (bool, error) {
	var exists bool
	var executor repository.SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `SELECT EXISTS(SELECT 1 FROM User WHERE email = ?)`
	err := executor.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *userRepository) GetUserByID(
	ctx context.Context,
	id string, opts ...repository.QueryOptions,
) (model.User, error) {
	var executor repository.SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM User
	WHERE id = ?
	`

	var user model.User
	err := executor.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
	)
	return user, err
}

func (ur *userRepository) GetUserByEmail(
	ctx context.Context,
	email string,
	opts ...repository.QueryOptions,
) (model.User, error) {
	var executor repository.SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM User
	WHERE email = ?
	`

	var user model.User
	err := executor.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
	)
	return user, err
}

func (ur *userRepository) Create(ctx context.Context, user *model.User, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = ur.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
    INSERT INTO User (
        id, name, email, password
    )
    VALUES (?, ?, ?, ?)
    `

	user.ID = uuid.New()

	_, err := executor.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)

	return err
}

func (ur *userRepository) Update(ctx context.Context, user model.User, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = ur.db
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
