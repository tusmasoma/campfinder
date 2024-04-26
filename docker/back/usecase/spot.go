//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type SpotUseCase interface {
	CreateSpot(ctx context.Context, params *CreateSpotParams) error
	BatchCreateSpots(ctx context.Context, params *BatchCreateSpotParams) error
	ListSpots(ctx context.Context, categories []string) []model.Spot
	GetSpot(ctx context.Context, spotID string) model.Spot
}

type spotUseCase struct {
	sr repository.SpotRepository
	cr repository.SpotsCacheRepository
}

func NewSpotUseCase(sr repository.SpotRepository, cr repository.SpotsCacheRepository) SpotUseCase {
	return &spotUseCase{
		sr: sr,
		cr: cr,
	}
}

type CreateSpotParams struct {
	Category    string
	Name        string
	Address     string
	Lat         float64
	Lng         float64
	Period      string
	Phone       string
	Price       string
	Description string
	IconPath    string
}

func (suc *spotUseCase) CreateSpot(ctx context.Context, params *CreateSpotParams) error {
	spots, err := suc.sr.List(
		ctx,
		[]repository.QueryCondition{
			{Field: "Lat", Value: params.Lat},
			{Field: "Lng", Value: params.Lng},
		},
	)
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return err
	}
	if len(spots) > 0 {
		log.Printf("Spot with this lat,lng already exists - status: %d", http.StatusConflict)
		return fmt.Errorf("already exists")
	}

	spot := model.Spot{
		Category:    params.Category,
		Name:        params.Name,
		Address:     params.Address,
		Lat:         params.Lat,
		Lng:         params.Lng,
		Period:      params.Period,
		Phone:       params.Phone,
		Price:       params.Price,
		Description: params.Description,
		IconPath:    params.IconPath,
	}

	if err = suc.sr.Create(ctx, spot); err != nil {
		log.Printf("Failed to create spot: %v", err)
		return err
	}
	return nil
}

type BatchCreateSpotParams struct {
	Spots []CreateSpotParams
}

func (suc *spotUseCase) BatchCreateSpots(ctx context.Context, params *BatchCreateSpotParams) error {
	var spots []model.Spot
	for _, param := range params.Spots {
		spot := model.Spot{
			ID:          uuid.New(),
			Category:    param.Category,
			Name:        param.Name,
			Address:     param.Address,
			Lat:         param.Lat,
			Lng:         param.Lng,
			Period:      param.Period,
			Phone:       param.Phone,
			Price:       param.Price,
			Description: param.Description,
			IconPath:    param.IconPath,
		}
		spots = append(spots, spot)
	}
	if err := suc.sr.BatchCreate(ctx, spots); err != nil {
		log.Printf("Failed to batch create spots: %v", err)
		return err
	}
	return nil
}

func (suc *spotUseCase) ListSpots(ctx context.Context, categories []string) []model.Spot {
	var allSpots []model.Spot
	var mu sync.Mutex
	wg := sync.WaitGroup{}

	for _, category := range categories {
		wg.Add(1)
		go func(category string) {
			defer wg.Done()

			spots, err := suc.sr.List(ctx, []repository.QueryCondition{{Field: "Category", Value: category}})
			if err != nil {
				log.Printf("Failed to get spot of %v: %v", category, err)
				spots = suc.getMasterData(ctx, category)
			} else {
				if cacheErr := suc.setMasterData(ctx, category, spots); cacheErr != nil {
					log.Printf("Failed to set master data of %v: %v", category, err)
				}
			}
			mu.Lock()
			allSpots = append(allSpots, spots...)
			mu.Unlock()
		}(category)
	}

	wg.Wait()
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
	spots, cacheErr := suc.cr.Get(ctx, "spots_"+category)
	if cacheErr != nil {
		log.Printf("Failed to get spots from cache for category %v: %v", category, cacheErr)
		return nil
	}
	return *spots
}

func (suc *spotUseCase) setMasterData(ctx context.Context, category string, spots []model.Spot) error {
	return suc.cr.Set(ctx, "spots_"+category, spots)
}
