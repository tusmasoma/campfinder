//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import "context"

type Base[T any] interface {
	List(ctx context.Context, qcs []QueryCondition) ([]T, error)
	Get(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, entity T) error
	BatchCreate(ctx context.Context, entitys []T) error
	Update(ctx context.Context, id string, entity T) error
	Delete(ctx context.Context, id string) error
}
