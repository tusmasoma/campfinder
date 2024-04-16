//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
)

type ImageRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]model.Image, error)
	Create(ctx context.Context, img model.Image) error
	Delete(ctx context.Context, id string) error
}

type ImagesCacheRepository interface {
	Set(ctx context.Context, key string, images model.Images) error
	Get(ctx context.Context, key string) (*model.Images, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Scan(ctx context.Context, match string) ([]string, error)
}
