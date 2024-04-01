package infra

import (
	"context"

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

func (ir *imageRepository) GetSpotImgURLBySpotID(
	ctx context.Context,
	spotID string,
	opts ...repository.QueryOptions,
) ([]model.Image, error) {
	var executor repository.SQLExecutor = ir.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM Image
	WHERE spot_id = ?
	`
	rows, err := executor.QueryContext(ctx, query, spotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var imgs []model.Image
	for rows.Next() {
		var img model.Image
		if err = rows.Scan(
			&img.ID,
			&img.SpotID,
			&img.UserID,
			&img.URL,
			&img.Created,
		); err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return imgs, nil
}
