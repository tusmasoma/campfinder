//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type ImageUseCase interface {
	ListImages(ctx context.Context, spotID string) ([]model.Image, error)
	CreateImage(ctx context.Context, spotID uuid.UUID, url string, user model.User) error
	DeleteImage(ctx context.Context, id string, userID string, user model.User) error
}

type imageUseCase struct {
	ir repository.ImageRepository
	ic repository.ImagesCacheRepository
}

func NewImageUseCase(ir repository.ImageRepository, ic repository.ImagesCacheRepository) ImageUseCase {
	return &imageUseCase{
		ir: ir,
		ic: ic,
	}
}

func (ih *imageUseCase) ListImages(ctx context.Context, spotID string) ([]model.Image, error) {
	images, err := ih.ir.List(ctx, []repository.QueryCondition{{Field: "SpotID", Value: spotID}})
	if err != nil {
		log.Printf("Failed to get images of %v: %v", spotID, err)
		images = ih.getMasterData(ctx, spotID)
		return images, nil
	}
	if cacheErr := ih.setMasterData(ctx, spotID, images); cacheErr != nil {
		log.Printf("Failed to set images data of %v: %v", spotID, cacheErr)
	}
	return images, nil
}

func (ih *imageUseCase) CreateImage(ctx context.Context, spotID uuid.UUID, url string, user model.User) error {
	img := model.Image{
		SpotID: spotID,
		UserID: user.ID,
		URL:    url,
	}
	if err := ih.ir.Create(ctx, img); err != nil {
		log.Printf("Failed to create image: %v", err)
		return err
	}
	return nil
}

func (ih *imageUseCase) DeleteImage(ctx context.Context, id string, userID string, user model.User) error {
	if !user.IsAdmin && user.ID.String() != userID {
		log.Print("Don't have permission to delete images")
		return fmt.Errorf("don't have permission to delete images")
	}

	if err := ih.ir.Delete(ctx, id); err != nil {
		log.Print("Internal server error while deleting image")
		return err
	}
	return nil
}

func (ih *imageUseCase) getMasterData(ctx context.Context, spotID string) []model.Image {
	images, cacheErr := ih.ic.Get(ctx, spotID)
	if cacheErr != nil {
		log.Printf("Failed to get images from cache for spotID %v: %v", spotID, cacheErr)
		return nil
	}
	return *images
}

func (ih *imageUseCase) setMasterData(ctx context.Context, spotID string, images []model.Image) error {
	return ih.ic.Set(ctx, spotID, images)
}
