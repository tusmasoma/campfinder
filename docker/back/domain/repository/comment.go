//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
)

type CommentRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]model.Comment, error)
	Get(ctx context.Context, id string) (*model.Comment, error)
	Create(ctx context.Context, comment model.Comment) error
	Update(ctx context.Context, id string, comment model.Comment) error
	Delete(ctx context.Context, id string) error
}

type CommentsCacheRepository interface {
	Set(ctx context.Context, key string, comments model.Comments) error
	Get(ctx context.Context, key string) (*model.Comments, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Scan(ctx context.Context, match string) ([]string, error)
}
