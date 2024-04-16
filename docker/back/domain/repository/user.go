//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
)

type UserRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user model.User) error
	Update(ctx context.Context, id string, spot model.User) error
	Delete(ctx context.Context, id string) error
}

type UserCacheRepository interface {
	Set(ctx context.Context, key string, user model.User) error
	SetUserSession(ctx context.Context, userID string, sessionData string) error
	Get(ctx context.Context, key string) (*model.User, error)
	GetUserSession(ctx context.Context, userID string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Scan(ctx context.Context, match string) ([]string, error)
}
