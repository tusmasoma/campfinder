package model

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID      uuid.UUID
	SpotID  uuid.UUID
	UserID  uuid.UUID
	URL     string
	Created time.Time
}
