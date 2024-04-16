//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package mysql

import (
	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type userRepository struct {
	*base[model.User]
}

func NewUserRepository(db repository.SQLExecutor, dialect *goqu.DialectWrapper) repository.UserRepository {
	return &userRepository{
		base: newBase[model.User](db, dialect, "User"),
	}
}
