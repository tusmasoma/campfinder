package model

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID      uuid.UUID `db:"id"`
	SpotID  uuid.UUID `db:"spot_id"`
	UserID  uuid.UUID `db:"user_id"`
	URL     string    `db:"url"`
	Created time.Time `db:"-"`
}

type Images []Image
