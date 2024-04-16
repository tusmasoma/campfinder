package mysql

import (
	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type imageRepository struct {
	*base[model.Image]
}

func NewImageRepository(db repository.SQLExecutor, dialect *goqu.DialectWrapper) repository.ImageRepository {
	return &imageRepository{
		base: newBase[model.Image](db, dialect, "Image"),
	}
}
