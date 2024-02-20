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

type CommentUseCase interface {
	GetCommentBySpotID(ctx context.Context, spotID string) ([]model.Comment, error)
	CommentCreate(ctx context.Context, spotID uuid.UUID, starRate float64, text string, user model.User) error
	CommentUpdate(
		ctx context.Context,
		id uuid.UUID,
		spotID uuid.UUID,
		userID uuid.UUID,
		starRate float64,
		text string,
		user model.User,
	) error
	CommentDelete(ctx context.Context, id string, userID string, user model.User) error
}

type commentUseCase struct {
	cr repository.CommentRepository
}

func NewCommentUseCase(cr repository.CommentRepository) CommentUseCase {
	return &commentUseCase{
		cr: cr,
	}
}

func (cuc *commentUseCase) GetCommentBySpotID(ctx context.Context, spotID string) ([]model.Comment, error) {
	return cuc.cr.GetCommentBySpotID(ctx, spotID)
}

func (cuc *commentUseCase) CommentCreate(
	ctx context.Context,
	spotID uuid.UUID,
	starRate float64,
	text string,
	user model.User,
) error {
	var comment = model.Comment{
		SpotID:   spotID,
		UserID:   user.ID,
		StarRate: starRate,
		Text:     text,
	}

	if err := cuc.cr.Create(ctx, comment); err != nil {
		log.Printf("Failed to create comment: %v", err)
		return err
	}
	return nil
}

func (cuc *commentUseCase) CommentUpdate(
	ctx context.Context,
	id uuid.UUID,
	spotID uuid.UUID,
	userID uuid.UUID,
	starRate float64,
	text string,
	user model.User,
) error {
	if !user.IsAdmin && user.ID != userID {
		log.Print("Don't have permission to update comment")
		return fmt.Errorf("don't have permission to update comment")
	}

	var comment = model.Comment{
		ID:       id,
		SpotID:   spotID,
		UserID:   userID,
		StarRate: starRate,
		Text:     text,
	}
	if err := cuc.cr.Update(ctx, comment); err != nil {
		log.Printf("Failed to update comment: %v", err)
		return err
	}
	return nil
}

func (cuc *commentUseCase) CommentDelete(ctx context.Context, id string, userID string, user model.User) error {
	if !user.IsAdmin && user.ID.String() != userID {
		log.Print("Don't have permission to delete comment")
		return fmt.Errorf("don't have permission to delete comment")
	}

	if err := cuc.cr.Delete(ctx, id); err != nil {
		log.Printf("Failed to delete comment: %v", err)
		return err
	}
	return nil
}
