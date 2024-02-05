//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/tusmasoma/campfinder/config"
	"github.com/tusmasoma/campfinder/domain/repository"
	"google.golang.org/api/option"
)

var ErrCacheMiss = errors.New("cache: key not found")

type AuthMiddleware interface {
	Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc
}

type authMiddleware struct {
	rr repository.CacheRepository
}

func NewAuthMiddleware(rr repository.CacheRepository) AuthMiddleware {
	return &authMiddleware{
		rr: rr,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (am *authMiddleware) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Firebase SDK のセットアップ
		opt := option.WithCredentialsFile(os.Getenv("CREDENTIALS"))
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error initializing Firebase app: %v\n", err), http.StatusInternalServerError)
			return
		}
		auth, err := app.Auth(context.Background())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error initializing Firebase Auth: %v\n", err), http.StatusInternalServerError)
			return
		}

		// リクエストヘッダにAuthorizationが存在するか確認
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authentication failed: missing Authorization header", http.StatusUnauthorized)
			return
		}

		// "Bearer "から始まるか確認
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Authorization failed: header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		jwt := parts[1]

		//　アクセストークンの検証
		log.Print(jwt)
		token, err := auth.VerifyIDToken(ctx, jwt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed 1: %v", err), http.StatusUnauthorized)
			return
		}

		// コンテキストに userID を保存
		ctx = context.WithValue(ctx, config.ContextUserIDKey, token.UID)

		nextFunc(w, r.WithContext(ctx))
	}
}
