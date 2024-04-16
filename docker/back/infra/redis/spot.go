package redis

import (
	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type spotsRepository struct {
	*base[model.Spots]
}

func NewSpotsRepository(client *redis.Client) repository.SpotsCacheRepository {
	return &spotsRepository{
		base: newBase[model.Spots](client),
	}
}
