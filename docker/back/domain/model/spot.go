package model

import "github.com/google/uuid"

type Spot struct {
	ID          uuid.UUID `db:"id"`
	Category    string    `db:"category" json:"category"`
	Name        string    `db:"name" json:"name"`
	Address     string    `db:"address" json:"address"`
	Lat         float64   `db:"lat" json:"lat"`
	Lng         float64   `db:"lng" json:"lng"`
	Period      string    `db:"period" json:"period"`
	Phone       string    `db:"phone" json:"phone"`
	Price       string    `db:"price" json:"price"`
	Description string    `db:"description" json:"description"`
	IconPath    string    `db:"iconpath" json:"iconpath"`
}
