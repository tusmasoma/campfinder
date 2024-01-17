//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ImageRepository interface {
	GetSpotImgURLBySpotID(ctx context.Context, spotID string, opts ...QueryOptions) (imgs []Image, err error)
	Create(ctx context.Context, img Image, opts ...QueryOptions) (err error)
	Delete(ctx context.Context, id string, opts ...QueryOptions) (err error)
}

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
	UserID  uuid.UUID
	URL     string
	Created time.Time
}

func (ir *imageRepository) GetSpotImgURLBySpotID(
	ctx context.Context,
	spotID string,
	opts ...QueryOptions,
) ([]Image, error) {
	var executor SQLExecutor = ir.db
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

	var imgs []Image
	for rows.Next() {
		var img Image
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

func (ir *imageRepository) Create(ctx context.Context, img Image, opts ...QueryOptions) error {
	var executor SQLExecutor = ir.db
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

func (ir *imageRepository) Delete(ctx context.Context, id string, opts ...QueryOptions) error {
	var executor SQLExecutor = ir.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}
	_, err := executor.ExecContext(ctx, "DELETE FROM Image WHERE id = ?", id)
	return err
}
