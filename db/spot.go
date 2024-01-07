package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type SpotRepository interface {
	CheckIfSpotExists(ctx context.Context, lat float64, lng float64, opts ...QueryOptions) (bool, error)
	GetSpotByCategory(ctx context.Context, category string, opts ...QueryOptions) (spots []Spot, err error)
	Create(ctx context.Context, spot Spot, opts ...QueryOptions) (err error)
	Update(ctx context.Context, spot Spot, opts ...QueryOptions) (err error)
	Delete(ctx context.Context, spot Spot, opts ...QueryOptions) (err error)
	UpdateOrCreate(ctx context.Context, spot Spot, opts ...QueryOptions) error
}

type spotRepository struct {
	db *sql.DB
}

func NewSpotRepository(db *sql.DB) SpotRepository {
	return &spotRepository{
		db: db,
	}
}

type Spot struct {
	ID          uuid.UUID
	Category    string  `json:"category"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Period      string  `json:"period"`
	Phone       string  `json:"phone"`
	Price       string  `json:"price"`
	Description string  `json:"description"`
	IconPath    string  `json:"iconpath"`
}

func (sr *spotRepository) CheckIfSpotExists(ctx context.Context, lat float64, lng float64, opts ...QueryOptions) (bool, error) {
	var exists bool
	var executor SQLExecutor = sr.db
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

func (sr *spotRepository) GetSpotByCategory(ctx context.Context, category string, opts ...QueryOptions) (spots []Spot, err error) {
	var executor SQLExecutor = sr.db
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
		return
	}
	defer rows.Close()

	for rows.Next() {
		var spot Spot
		err = rows.Scan(
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
		spots = append(spots, spot)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (sr *spotRepository) Create(ctx context.Context, spot Spot, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	    INSERT INTO spot (
		category, name, address, lat, lng,
		period, phone, price, description, iconpath
		)
		VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
		`

	err = executor.QueryRowContext(
		ctx,
		query,
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
	).Scan(&spot.ID)

	return
}

func (sr *spotRepository) Update(ctx context.Context, spot Spot, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	UPDATE Spot SET
	name=?, period=?, phone=?, price=?, description=?
	WHERE id = ?
	`
	_, err = executor.ExecContext(ctx, query, spot.Name, spot.Period, spot.Phone, spot.Price, spot.Description, spot.ID)
	return
}

func (sr *spotRepository) Delete(ctx context.Context, spot Spot, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = sr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	_, err = executor.ExecContext(ctx, "DELETE FROM Spot WHERE id = ?", spot.ID)
	return
}

func (sr *spotRepository) UpdateOrCreate(ctx context.Context, spot Spot, opts ...QueryOptions) error {

	exists, err := sr.CheckIfSpotExists(ctx, spot.Lat, spot.Lng, opts...)
	if err != nil {
		return err
	}

	if exists {
		return sr.Update(ctx, spot, opts...)
	}
	return sr.Create(ctx, spot, opts...)
}
