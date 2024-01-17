//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tusmasoma/campfinder/cache"
	"github.com/tusmasoma/campfinder/db"
	"github.com/tusmasoma/campfinder/internal/auth"
)

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler interface {
	HandleUserCreate(w http.ResponseWriter, r *http.Request)
	HandleUserLogin(w http.ResponseWriter, r *http.Request)
	HandleUserLogout(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	ur db.UserRepository
	rr cache.RedisRepository
	ah auth.Handler
}

func NewUserHandler(ur db.UserRepository, rr cache.RedisRepository, ah auth.Handler) UserHandler {
	return &userHandler{
		ur: ur,
		rr: rr,
		ah: ah,
	}
}

func (uh *userHandler) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// リクエストの検証
	var requestBody UserCreateRequest
	if ok := isValidUserCreateRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// ユーザー登録
	user, err := uh.CreateUser(ctx, requestBody)
	if err != nil {
		http.Error(w, "Internal server error while creating user", http.StatusInternalServerError)
		return
	}

	// アクセストークンの生成とキャッシュへの保存
	jwt, err := uh.GenerateAndStoreToken(ctx, user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ヘッダーにアクセストークンをセット
	w.Header().Set("Authorization", "Bearer "+jwt)

	// レスポンスボディやその他のヘッダーをセットして、レスポンスを送信する
	w.WriteHeader(http.StatusOK)
}

func isValidUserCreateRequest(body io.ReadCloser, requestBody *UserCreateRequest) bool {
	// リクエストボディのJSONを構造体にデコード
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Printf("Missing required fields: Name or Password")
		return false
	}
	return true
}

func (uh *userHandler) CreateUser(ctx context.Context, requestBody UserCreateRequest) (db.User, error) {
	// usernameが登録済みかどうかMySQLに問い合わせ
	exists, err := uh.ur.CheckIfUserExists(ctx, requestBody.Email)
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return db.User{}, err
	}
	if exists {
		// ユーザーがすでに存在している場合、409 Conflictを返す
		log.Printf("User with this name already exists - status: %d", http.StatusConflict)
		return db.User{}, fmt.Errorf("user with this name already exists")
	}

	// ユーザ情報登録
	var user db.User
	user.Email = requestBody.Email
	user.Name = ExtractUsernameFromEmail(requestBody.Email)
	password, err := db.PasswordEncrypt(requestBody.Password)
	if err != nil { // paswordのハッシュ化
		log.Printf("Internal server error: %v", err)
		return db.User{}, err
	}
	user.Password = password

	if err = uh.ur.Create(ctx, &user); err != nil {
		log.Printf("Failed to create user: %v", err)
		return db.User{}, err
	}
	return user, nil
}

func ExtractUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func (uh *userHandler) GenerateAndStoreToken(ctx context.Context, user db.User) (string, error) {
	// アクセストークンを生成
	jwt, jti := auth.GenerateToken(user)

	// Cacheに保存
	if err := uh.rr.Set(ctx, user.ID.String(), jti); err != nil {
		log.Print("Failed to set access token in cache")
		return "", err
	}
	return jwt, nil
}

func (uh *userHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// リクエストの検証
	var requestBody UserLoginRequest
	if ok := isValidUserLoginRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// emailでMySQLにユーザー情報問い合わせ
	user, err := uh.ur.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		http.Error(w, "Error retrieving user by email", http.StatusInternalServerError)
		return
	}

	// 既にログイン済みかどうか確認する
	isAuthenticate := uh.rr.Exists(ctx, user.ID.String())
	if isAuthenticate {
		http.Error(w, "Already logged in", http.StatusInternalServerError)
		return
	}

	// Clientから送られてきたpasswordをハッシュ化したものとMySQLから返されたハッシュ化されたpasswordを比較する
	if err = db.CompareHashAndPassword(user.Password, requestBody.Password); err != nil {
		http.Error(w, "password does not match", http.StatusInternalServerError)
		return
	}

	// アクセストークンの生成とキャッシュへの保存
	jwt, err := uh.GenerateAndStoreToken(ctx, user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)

	w.WriteHeader(http.StatusOK)
}

func isValidUserLoginRequest(body io.ReadCloser, requestBody *UserLoginRequest) bool {
	// リクエストボディのJSONを構造体にデコード
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Printf("Missing required fields: Name or Password")
		return false
	}
	return true
}

func (uh *userHandler) HandleUserLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ユーザー情報取得
	user, err := uh.ah.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	// Redisに該当のuserIDの削除問い合わせ
	if err = uh.rr.Delete(ctx, user.ID.String()); err != nil {
		http.Error(w, "Failed to delete userID from redis", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
