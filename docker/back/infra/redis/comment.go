package redis

import (
	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type commentsRepository struct {
	*base[model.Comments]
}

func NewCommentsRepository(client *redis.Client) repository.CommentsCacheRepository {
	return &commentsRepository{
		base: newBase[model.Comments](client),
	}
}
