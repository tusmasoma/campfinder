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
	cr repository.CacheRepository
}

func NewSpotUseCase(sr repository.SpotRepository, cr repository.CacheRepository) SpotUseCase {
	return &spotUseCase{
		sr: sr,
		cr: cr,
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

	for _, category := range categories {
		spots, err := suc.sr.List(ctx, []repository.QueryCondition{{Field: "Category", Value: category}})
		if err != nil {
			log.Printf("Failed to get spot of %v: %v", category, err)
			spots = suc.getMasterData(ctx, category)
			allSpots = append(allSpots, spots...)
			continue
		}
		if cacheErr := suc.setMasterData(ctx, category, spots); cacheErr != nil {
			log.Printf("Failed to set master data of %v: %v", category, err)
		}
		allSpots = append(allSpots, spots...)
	}

	return allSpots
}

func (suc *spotUseCase) GetSpot(ctx context.Context, spotID string) model.Spot {
	spot, err := suc.sr.Get(ctx, spotID)
	if err != nil {
		log.Printf("Failed to get spot of %v: %v", spotID, err)

		var allSpots []model.Spot
		keys, scanErr := suc.cr.Scan(ctx, "spots_*")
		if scanErr != nil {
			log.Printf("Failed to scan cache: %v", scanErr)
			return model.Spot{}
		}
		for _, key := range keys {
			category := strings.TrimPrefix(key, "spots_")
			spots := suc.getMasterData(ctx, category)
			allSpots = append(allSpots, spots...)
		}
		for _, spot := range allSpots {
			if spot.ID.String() == spotID {
				return spot
			}
		}
		return model.Spot{}
	}
	return *spot
}

func (suc *spotUseCase) getMasterData(ctx context.Context, category string) []model.Spot {
	temp, cacheErr := suc.cr.Get(ctx, "spots_"+category)
	if cacheErr != nil {
		log.Printf("Failed to get spots from cache for category %v: %v", category, cacheErr)
		return nil
	}
	spots, ok := temp.([]model.Spot)
	if !ok {
		return nil
	}
	return spots
}

func (suc *spotUseCase) setMasterData(ctx context.Context, category string, spots []model.Spot) error {
	return suc.cr.Set(ctx, "spots_"+category, spots)
}
