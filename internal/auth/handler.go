//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/tusmasoma/campfinder/db"
)

type ContextKey string

const ContextUserIDKey ContextKey = "userID"

type Handler interface {
	FetchUserFromContext(ctx context.Context) (db.User, error)
}

type handler struct {
	ur db.UserRepository
}

func NewAuthHandler(ur db.UserRepository) Handler {
	return &handler{
		ur: ur,
	}
}

func (ah *handler) FetchUserFromContext(ctx context.Context) (db.User, error) {
	// リクエストのコンテキストからユーザー名を取得し、そのユーザー名を用いてデータベースからユーザー情報を返す
	userID, ok := ctx.Value(ContextUserIDKey).(string)
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
