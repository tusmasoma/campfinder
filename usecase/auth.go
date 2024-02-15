//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/tusmasoma/campfinder/config"
	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository"
)

type AuthUseCase interface {
	FetchUserFromContext(ctx context.Context) (*model.User, error)
}

type authUseCase struct {
	ur repository.UserRepository
}

func NewAuthUseCase(ur repository.UserRepository) AuthUseCase {
	return &authUseCase{
		ur: ur,
	}
}

func (auc *authUseCase) FetchUserFromContext(ctx context.Context) (*model.User, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Printf("Failed to retrieve userId from context")
		return nil, fmt.Errorf("user name not found in request context")
	}
	user, err := auc.ur.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Failed to get UserInfo from db: %v", userID)
		return nil, err
	}
	return &user, nil
}
