//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/campfinder/domain/model"
)

type SpotRepository interface {
	CheckIfSpotExists(ctx context.Context, lat float64, lng float64, opts ...QueryOptions) (bool, error)
	GetSpotByID(ctx context.Context, id string, opts ...QueryOptions) (model.Spot, error)
	GetSpotByCategory(ctx context.Context, category string, opts ...QueryOptions) (spots []model.Spot, err error)
	Create(ctx context.Context, spot model.Spot, opts ...QueryOptions) (err error)
	Update(ctx context.Context, spot model.Spot, opts ...QueryOptions) (err error)
	Delete(ctx context.Context, id string, opts ...QueryOptions) (err error)
	UpdateOrCreate(ctx context.Context, spot model.Spot, opts ...QueryOptions) error
}
