//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package infra

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type spotRepository struct {
	db *sql.DB
}

func NewSpotRepository(db *sql.DB) repository.SpotRepository {
	return &spotRepository{
		db: db,
	}
}

func (sr *spotRepository) CheckIfSpotExists(
	ctx context.Context,
	lat float64,
	lng float64,
	opts ...repository.QueryOptions,
) (bool, error) {
	var exists bool
	var executor repository.SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `SELECT exists(SELECT 1 FROM Spot WHERE lat = ? AND lng = ?)`
	err := executor.QueryRowContext(ctx, query, lat, lng).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (sr *spotRepository) GetSpotByID(
	ctx context.Context,
	id string,
	opts ...repository.QueryOptions,
) (model.Spot, error) {
	var executor repository.SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM Spot
	WHERE id = ?
	`

	var spot model.Spot
	err := executor.QueryRowContext(ctx, query, id).Scan(
		&spot.ID,
		&spot.Category,
		&spot.Name,
		&spot.Address,
		&spot.Lat,
		&spot.Lng,
		&spot.Period,
		&spot.Phone,
		&spot.Price,
		&spot.Description,
		&spot.IconPath,
	)
	return spot, err
}

func (sr *spotRepository) GetSpotByCategory(
	ctx context.Context,
	category string,
	opts ...repository.QueryOptions,
) ([]model.Spot, error) {
	var executor repository.SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT id, category, name, address, lat, lng, period, phone, price, description, iconpath
	FROM Spot
	WHERE category = ?
	`
	rows, err := executor.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spots []model.Spot
	for rows.Next() {
		var spot model.Spot
		if err = rows.Scan(
			&spot.ID,
			&spot.Category,
			&spot.Name,
			&spot.Address,
			&spot.Lat,
			&spot.Lng,
			&spot.Period,
			&spot.Phone,
			&spot.Price,
			&spot.Description,
			&spot.IconPath,
		); err != nil {
			return nil, err
		}
		spots = append(spots, spot)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return spots, nil
}

func (sr *spotRepository) Create(ctx context.Context, spot model.Spot, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	    INSERT INTO Spot (
		id, category, name, address, lat, lng,
		period, phone, price, description, iconpath
		)
		VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

	_, err := executor.ExecContext(
		ctx,
		query,
		uuid.New(),
		spot.Category,
		spot.Name,
		spot.Address,
		spot.Lat,
		spot.Lng,
		spot.Period,
		spot.Phone,
		spot.Price,
		spot.Description,
		spot.IconPath,
	)

	return err
}

func (sr *spotRepository) Update(ctx context.Context, spot model.Spot, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	UPDATE Spot SET
	name=?, period=?, phone=?, price=?, description=?
	WHERE id = ?
	`
	_, err := executor.ExecContext(ctx, query, spot.Name, spot.Period, spot.Phone, spot.Price, spot.Description, spot.ID)
	return err
}

func (sr *spotRepository) Delete(ctx context.Context, id string, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	_, err := executor.ExecContext(ctx, "DELETE FROM Spot WHERE id = ?", id)
	return err
}

func (sr *spotRepository) UpdateOrCreate(ctx context.Context, spot model.Spot, opts ...repository.QueryOptions) error {
	exists, err := sr.CheckIfSpotExists(ctx, spot.Lat, spot.Lng, opts...)
	if err != nil {
		return err
	}

	if exists {
		return sr.Update(ctx, spot, opts...)
	}
	return sr.Create(ctx, spot, opts...)
}
