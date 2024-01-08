//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/tusmasoma/campfinder/db"
)

type AuthHandler interface {
	FetchUserFromContext(ctx context.Context) (db.User, error)
}

type authHandler struct {
	ur db.UserRepository
}

func NewAuthHandler(ur db.UserRepository) AuthHandler {
	return &authHandler{
		ur: ur,
	}
}

func (ah *authHandler) FetchUserFromContext(ctx context.Context) (db.User, error) {
	// リクエストのコンテキストからユーザー名を取得し、そのユーザー名を用いてデータベースからユーザー情報を返す
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		log.Printf("Failed to retrieve userId from context")
		return db.User{}, fmt.Errorf("user name not found in request context")
	}
	user, err := ah.ur.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Failed to get UserInfo from db: %v", userID)
		return db.User{}, err
	}
	return user, nil
}
