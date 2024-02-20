package model

import "github.com/google/uuid"

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
