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
	ListComments(ctx context.Context, spotID string) ([]model.Comment, error)
	CreateComment(ctx context.Context, spotID uuid.UUID, starRate float64, text string, user model.User) error
	BatchCreateComments(ctx context.Context, params *BatchCreateCommentsParams) error
	UpdateComment(
		ctx context.Context,
		id uuid.UUID,
		spotID uuid.UUID,
		userID uuid.UUID,
		starRate float64,
		text string,
		user model.User,
	) error
	DeleteComment(ctx context.Context, id string, userID string, user model.User) error
}

type commentUseCase struct {
	cr repository.CommentRepository
	cc repository.CommentsCacheRepository
}

func NewCommentUseCase(cr repository.CommentRepository, cc repository.CommentsCacheRepository) CommentUseCase {
	return &commentUseCase{
		cr: cr,
		cc: cc,
	}
}

func (cuc *commentUseCase) ListComments(ctx context.Context, spotID string) ([]model.Comment, error) {
	comments, err := cuc.cr.List(ctx, []repository.QueryCondition{{Field: "SpotID", Value: spotID}})
	if err != nil {
		log.Printf("Failed to get comments of %v: %v", spotID, err)
		comments = cuc.getMasterData(ctx, spotID)
		return comments, nil
	}
	if cacheErr := cuc.setMasterData(ctx, spotID, comments); cacheErr != nil {
		log.Printf("Failed to set comments data of %v: %v", spotID, cacheErr)
	}
	return comments, nil
}

type CreateCommentParams struct {
	SpotID   uuid.UUID
	StarRate float64
	Text     string
	userID   uuid.UUID
}

func (cuc *commentUseCase) CreateComment(
	ctx context.Context,
	spotID uuid.UUID,
	starRate float64,
	text string,
	user model.User,
) error {
	comment := model.Comment{
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

type BatchCreateCommentsParams struct {
	Comments []CreateCommentParams
}

func (cuc *commentUseCase) BatchCreateComments(ctx context.Context, params *BatchCreateCommentsParams) error {
	var comments []model.Comment
	for _, param := range params.Comments {
		comment := model.Comment{
			SpotID:   param.SpotID,
			UserID:   param.userID,
			StarRate: param.StarRate,
			Text:     param.Text,
		}
		comments = append(comments, comment)
	}
	if err := cuc.cr.BatchCreate(ctx, comments); err != nil {
		log.Printf("Failed to batch create comments: %v", err)
		return err
	}
	return nil
}

func (cuc *commentUseCase) UpdateComment(
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

	comment := model.Comment{
		ID:       id,
		SpotID:   spotID,
		UserID:   userID,
		StarRate: starRate,
		Text:     text,
	}
	if err := cuc.cr.Update(ctx, id.String(), comment); err != nil {
		log.Printf("Failed to update comment: %v", err)
		return err
	}
	return nil
}

func (cuc *commentUseCase) DeleteComment(ctx context.Context, id string, userID string, user model.User) error {
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

func (cuc *commentUseCase) getMasterData(ctx context.Context, spotID string) []model.Comment {
	comments, cacheErr := cuc.cc.Get(ctx, "comments_"+spotID)
	if cacheErr != nil {
		log.Printf("Failed to get comments from cache for spotID %v: %v", spotID, cacheErr)
		return nil
	}
	return *comments
}

func (cuc *commentUseCase) setMasterData(ctx context.Context, spotID string, comments []model.Comment) error {
	return cuc.cc.Set(ctx, "comments_"+spotID, comments)
}
