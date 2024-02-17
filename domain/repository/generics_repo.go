package repository

import "context"

type Repository[T any] interface {
	GetByID(ctx context.Context, id string) (T, error)
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id string) error
}
