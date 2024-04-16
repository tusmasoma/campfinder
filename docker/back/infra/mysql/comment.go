//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package mysql

import (
	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type commentRepository struct {
	*base[model.Comment]
}

func NewCommentRepository(db repository.SQLExecutor, dialect *goqu.DialectWrapper) repository.CommentRepository {
	return &commentRepository{
		base: newBase[model.Comment](db, dialect, "Comment"),
	}
}
