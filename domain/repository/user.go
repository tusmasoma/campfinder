package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/domain/model"
)

type UserRepository interface {
	CheckIfUserExists(ctx context.Context, email string, opts ...QueryOptions) (bool, error)
	GetUserByID(ctx context.Context, id string, opts ...QueryOptions) (model.User, error)
	GetUserByEmail(ctx context.Context, name string, opts ...QueryOptions) (model.User, error)
	Create(ctx context.Context, user *model.User, opts ...QueryOptions) error
	Update(ctx context.Context, user model.User, opts ...QueryOptions) error
}
