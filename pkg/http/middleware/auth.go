package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tusmasoma/campfinder/cache"
	"github.com/tusmasoma/campfinder/internal/auth"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type AuthMiddleware interface {
	Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc
}

type authMiddleware struct {
	rr cache.RedisRepository
}

func NewAuthMiddleware(rr cache.RedisRepository) AuthMiddleware {
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
		err := auth.ValidateAccessToken(jwt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		// JWTからペイロード取得
		var payload auth.Payload
		payload, err = auth.GetPayloadFromToken(jwt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		// 該当のuserIdが存在するかキャッシュに問い合わせ
		jti, err := am.rr.Get(ctx, payload.UserID)
		if errors.Is(err, cache.ErrCacheMiss) {
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
		ctx = context.WithValue(ctx, auth.ContextUserIDKey, payload.UserID)

		nextFunc(w, r.WithContext(ctx))
	}
}

// 以下のミドルウェアは、検証okのtokenをコンテキスに渡すので、tokenから情報を取得する処理は別でやらないといけない
// middleware/jwt.go

const cacheDuration = 5 * time.Minute

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(_ context.Context) error {
	return nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, cacheDuration)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte(`{"message":"Failed to validate JWT."}`))
		if err != nil {
			log.Printf("エラー: %v", err)
		}
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}
