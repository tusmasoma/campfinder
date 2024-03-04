package repository

import "context"

type Repository[T any] interface {
	List(ctx context.Context, qcs []QueryCondition) ([]T, error)
	Get(ctx context.Context, id string) (T, error)
	Create(ctx context.Context, entity []T) error
	Update(ctx context.Context, id string, entity T) error
	Delete(ctx context.Context, id string) error
}
