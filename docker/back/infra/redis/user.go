package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type userRepository struct {
	*base[model.User]
}

func NewUserRepository(client *redis.Client) repository.UserCacheRepository {
	return &userRepository{
		base: newBase[model.User](client),
	}
}

func (ur *userRepository) GetUserSession(ctx context.Context, userID string) (string, error) {
	return ur.client.Get(ctx, userID).Result()
}

func (ur *userRepository) SetUserSession(ctx context.Context, userID string, sessionData string) error {
	return ur.client.Set(ctx, userID, sessionData, 0).Err()
}
