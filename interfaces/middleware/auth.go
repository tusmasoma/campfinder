//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tusmasoma/campfinder/config"
	"github.com/tusmasoma/campfinder/domain/repository"
	"github.com/tusmasoma/campfinder/internal/auth"
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
		err := auth.ValidateAccessToken(jwt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed 1: %v", err), http.StatusUnauthorized)
			return
		}

		// JWTからペイロード取得
		var payload auth.Payload
		payload, err = auth.GetPayloadFromToken(jwt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed 2: %v", err), http.StatusUnauthorized)
			return
		}

		// 該当のuserIdが存在するかキャッシュに問い合わせ
		jti, err := am.rr.Get(ctx, payload.UserID)
		if errors.Is(err, ErrCacheMiss) {
			http.Error(w, "Authentication failed: userId is not exit on cache", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Authentication failed: missing userId on cache", http.StatusUnauthorized)
			return
		}

		// Redisから取得したjtiとJWTのjtiを比較
		if payload.JTI != jti {
			http.Error(w, "Authentication failed: jwt does not match", http.StatusUnauthorized)
			return
		}

		// 今後有効期限の確認も行う

		// コンテキストに userID を保存
		ctx = context.WithValue(ctx, config.ContextUserIDKey, payload.UserID)

		nextFunc(w, r.WithContext(ctx))
	}
}