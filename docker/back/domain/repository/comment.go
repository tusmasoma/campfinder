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
