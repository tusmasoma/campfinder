//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
)

type SpotRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]model.Spot, error)
	Get(ctx context.Context, id string) (*model.Spot, error)
	Create(ctx context.Context, spot model.Spot) error
	Update(ctx context.Context, id string, spot model.Spot) error
	Delete(ctx context.Context, id string) error
	CreateOrUpdate(ctx context.Context, id string, qcs []QueryCondition, spot model.Spot) error
}

type SpotsCacheRepository interface {
	Set(ctx context.Context, key string, spots model.Spots) error
	Get(ctx context.Context, key string) (*model.Spots, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Scan(ctx context.Context, match string) ([]string, error)
}
