package redis

import (
	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type imagesRepository struct {
	*base[model.Images]
}

func NewImagesRepository(client *redis.Client) repository.ImagesCacheRepository {
	return &imagesRepository{
		base: newBase[model.Images](client),
	}
}
