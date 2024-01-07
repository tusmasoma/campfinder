package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ImageRepository interface{}

type imageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) ImageRepository {
	return &imageRepository{
		db: db,
	}
}

type Image struct {
	ID      uuid.UUID
	SpotID  uuid.UUID
	URL     string
	Created time.Time
}

func (ir *imageRepository) GetSpotImgURLBySpotId(ctx context.Context, spotId uuid.UUID, opts ...QueryOptions) (imgs []Image, err error) {
	var executor SQLExecutor = ir.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM Image
	WHERE spot_id = ?
	`
	rows, err := executor.QueryContext(ctx, query, spotId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var img Image
		err = rows.Scan(
			&img.ID,
			&img.SpotID,
			&img.URL,
			&img.Created,
		)
		imgs = append(imgs, img)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (ir *imageRepository) Create(ctx context.Context, img Image, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = ir.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	INSERT INTO Image (
		spot_id, url
		)
		VALUES (?, ?)
		RETURNING id;
		`
	err = executor.QueryRowContext(
		ctx,
		query,
		img.SpotID,
		img.URL,
	).Scan(&img.ID)

	return
}

func (ir *imageRepository) Delete(ctx context.Context, id uuid.UUID, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = ir.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}
	_, err = executor.ExecContext(ctx, "DELETE FROM Image WHERE id = ?", id)
	return
}
