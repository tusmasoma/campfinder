//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
)

type ImageRepository interface {
	GetSpotImgURLBySpotID(ctx context.Context, spotID string, opts ...QueryOptions) (imgs []model.Image, err error)
	List(ctx context.Context, qcs []QueryCondition) ([]model.Image, error)
	Create(ctx context.Context, img model.Image) error
	Delete(ctx context.Context, id string) error
}
