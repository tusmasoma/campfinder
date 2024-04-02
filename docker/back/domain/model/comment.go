package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID       uuid.UUID `db:"id"`
	SpotID   uuid.UUID `db:"spot_id"`
	UserID   uuid.UUID `db:"user_id"`
	StarRate float64   `db:"star_rate" json:"starRate"`
	Text     string    `db:"text" json:"text"`
	Created  time.Time `db:"-"`
}
