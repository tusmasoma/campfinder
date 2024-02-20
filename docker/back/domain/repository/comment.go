//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
)

type CommentRepository interface {
	GetCommentBySpotID(ctx context.Context, spotID string, opts ...QueryOptions) (comments []model.Comment, err error)
	GetCommentByID(ctx context.Context, id string, opts ...QueryOptions) (comment model.Comment, err error)
	Create(ctx context.Context, comment model.Comment, opts ...QueryOptions) (err error)
	Update(ctx context.Context, comment model.Comment, opts ...QueryOptions) (err error)
	Delete(ctx context.Context, id string, opts ...QueryOptions) (err error)
}
