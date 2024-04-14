//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package infra

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

var ErrCacheMiss = errors.New("cache: key not found")

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) repository.CacheRepository {
	return &redisRepository{
		client: client,
	}
}

func (rr *redisRepository) Set(ctx context.Context, key string, value interface{}) error {
	err := rr.client.Set(ctx, key, value, 0).Err()
	return err
}

func (rr *redisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := rr.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCacheMiss
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (rr *redisRepository) Delete(ctx context.Context, key string) error {
	err := rr.client.Del(ctx, key).Err()
	return err
}

func (rr *redisRepository) Exists(ctx context.Context, key string) bool {
	val := rr.client.Exists(ctx, key).Val()
	return val > 0
}

func (rr *redisRepository) Scan(ctx context.Context, match string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	for {
		keys, newCursor, err := rr.client.Scan(ctx, cursor, match, 0).Result()
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	return allKeys, nil
}
