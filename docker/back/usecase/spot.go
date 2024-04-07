//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type SpotUseCase interface {
	CreateSpot(
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
	ListSpots(ctx context.Context, categories []string) []model.Spot
	GetSpot(ctx context.Context, spotID string) model.Spot
}

type spotUseCase struct {
	sr repository.SpotRepository
	rr repository.CacheRepository
}

func NewSpotUseCase(sr repository.SpotRepository) SpotUseCase {
	return &spotUseCase{
		sr: sr,
	}
}

func (suc *spotUseCase) CreateSpot(
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

func (suc *spotUseCase) ListSpots(ctx context.Context, categories []string) []model.Spot {
	var allSpots []model.Spot
	var err error

	for _, category := range categories {
		spots, err := suc.sr.List(ctx, []repository.QueryCondition{{Field: "Category", Value: category}})
		if err != nil {
			log.Printf("Failed to get spot of %v: %v", category, err)
			continue
		}
		if cacheErr := suc.setMasterData(ctx, category, spots); cacheErr != nil {
			log.Printf("Failed to set master data of %v: %v", category, err)
			continue
		}
		allSpots = append(allSpots, spots...)
	}

	if err != nil {
		allSpots = suc.getMasterData(ctx, categories)
		return allSpots
	}

	return allSpots
}

func (suc *spotUseCase) GetSpot(ctx context.Context, spotID string) model.Spot {
	spot, err := suc.sr.Get(ctx, spotID)
	if err != nil {
		log.Printf("Failed to get spot of %v: %v", spotID, err)

		var categories []string
		keys, scanErr := suc.rr.Scan(ctx, "spots_*")
		if scanErr != nil {
			log.Printf("Failed to scan cache: %v", scanErr)
			return model.Spot{}
		}
		for _, key := range keys {
			category := strings.TrimPrefix(key, "spots_")
			categories = append(categories, category)
		}
		spots := suc.getMasterData(ctx, categories)
		for _, spot := range spots {
			if spot.ID.String() == spotID {
				return spot
			}
		}
	}
	return *spot
}

func (suc *spotUseCase) getMasterData(ctx context.Context, categories []string) []model.Spot {
	var allSpots []model.Spot
	for _, category := range categories {
		temp, cacheErr := suc.rr.Get(ctx, "spots_"+category)
		if cacheErr != nil {
			log.Printf("Failed to get spots from cache for category %v: %v", category, cacheErr)
			continue
		}
		allSpots = append(allSpots, temp.([]model.Spot)...)
	}
	return allSpots
}

func (suc *spotUseCase) setMasterData(ctx context.Context, category string, spots []model.Spot) error {
	return suc.rr.Set(ctx, "spots_"+category, spots)
}
