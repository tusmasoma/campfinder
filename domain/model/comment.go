package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID       uuid.UUID
	SpotID   uuid.UUID
	UserID   uuid.UUID
	StarRate float64 `json:"starRate"`
	Text     string  `json:"text"`
	Created  time.Time
}
