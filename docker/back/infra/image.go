package infra

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type imageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) repository.ImageRepository {
	return &imageRepository{
		db: db,
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

func (ir *imageRepository) Create(ctx context.Context, img model.Image, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = ir.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	INSERT INTO Image (
		id, spot_id, user_id, url
		)
		VALUES (?, ?, ?, ?)
		`
	_, err := executor.ExecContext(
		ctx,
		query,
		uuid.New(),
		img.SpotID,
		img.UserID,
		img.URL,
	)

	return err
}

func (ir *imageRepository) Delete(ctx context.Context, id string, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = ir.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}
	_, err := executor.ExecContext(ctx, "DELETE FROM Image WHERE id = ?", id)
	return err
}
