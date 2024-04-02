//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type SpotUseCase interface {
	SpotCreate(
		ctx context.Context,
		category string,
		name string,
		address string,
		lat float64,
		lng float64,
		period string,
		phone string,
		price string,
		description string,
		iconPath string,
	) error
	SpotGet(ctx context.Context, categories []string, spotID string) []model.Spot
}

type spotUseCase struct {
	sr repository.SpotRepository
}

func NewSpotUseCase(sr repository.SpotRepository) SpotUseCase {
	return &spotUseCase{
		sr: sr,
	}
}

func (suc *spotUseCase) SpotCreate(
	ctx context.Context,
	category string,
	name string,
	address string,
	lat float64,
	lng float64,
	period string,
	phone string,
	price string,
	description string,
	iconPath string,
) error {
	spots, err := suc.sr.List(ctx, []repository.QueryCondition{{Field: "Lat", Value: lat}, {Field: "Lng", Value: lng}})
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return err
	}
	if len(spots) > 0 {
		log.Printf("Spot with this lat,lng already exists - status: %d", http.StatusConflict)
		return fmt.Errorf("already exists")
	}

	spot := model.Spot{
		Category:    category,
		Name:        name,
		Address:     address,
		Lat:         lat,
		Lng:         lng,
		Period:      period,
		Phone:       phone,
		Price:       price,
		Description: description,
		IconPath:    iconPath,
	}

	if err = suc.sr.Create(ctx, spot); err != nil {
		log.Printf("Failed to create spot: %v", err)
		return err
	}
	return nil
}

func (suc *spotUseCase) SpotGet(ctx context.Context, categories []string, spotID string) []model.Spot {
	var allSpots []model.Spot

	for _, category := range categories {
		spots, err := suc.sr.List(ctx, []repository.QueryCondition{{Field: "Category", Value: category}})
		if err != nil {
			log.Printf("Failed to get spot of %v: %v", category, err)
			continue
		}
		allSpots = append(allSpots, spots...)
	}

	if spotID != "" {
		spot, err := suc.sr.Get(ctx, spotID)
		if err != nil {
			log.Printf("Failed to get spot of %v: %v", spotID, err)
		}
		allSpots = append(allSpots, *spot)
	}

	return allSpots
}
