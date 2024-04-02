//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package infra

import (
	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type spotRepository struct {
	*base[model.Spot]
}

func NewSpotRepository(db repository.SQLExecutor, dialect *goqu.DialectWrapper) repository.SpotRepository {
	return &spotRepository{
		base: newBase[model.Spot](db, dialect, "Spot"),
	}
}
