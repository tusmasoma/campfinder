//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository"
)

type ImageUseCase interface {
	GetSpotImgURLBySpotID(ctx context.Context, spotID string) ([]model.Image, error)
	ImageCreate(ctx context.Context, spotID uuid.UUID, url string, user model.User) error
	ImageDelete(ctx context.Context, id string, userID string, user model.User) error
}

type imageUseCase struct {
	ir repository.ImageRepository
}

func NewImageUseCase(ir repository.ImageRepository) ImageUseCase {
	return &imageUseCase{
		ir: ir,
	}
}

func (ih *imageUseCase) GetSpotImgURLBySpotID(ctx context.Context, spotID string) ([]model.Image, error) {
	return ih.ir.GetSpotImgURLBySpotID(ctx, spotID)
}

func (ih *imageUseCase) ImageCreate(ctx context.Context, spotID uuid.UUID, url string, user model.User) error {
	var img = model.Image{
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

func (ih *imageUseCase) ImageDelete(ctx context.Context, id string, userID string, user model.User) error {
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
