//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/tusmasoma/campfinder/internal/auth"
	"github.com/tusmasoma/campfinder/usecase"
)

type UserHandler interface {
	HandleUserCreate(w http.ResponseWriter, r *http.Request)
	HandleUserLogin(w http.ResponseWriter, r *http.Request)
	HandleUserLogout(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	uur usecase.UserUseCase
	ah  auth.Handler
}

func NewUserHandler(uur usecase.UserUseCase, ah auth.Handler) UserHandler {
	return &userHandler{
		uur: uur,
		ah:  ah,
	}
}

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody UserCreateRequest
	if ok := isValidUserCreateRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jwt, err := uh.uur.CreateUserAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)
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

func (uh *userHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody UserLoginRequest
	if ok := isValidUserLoginRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jwt, err := uh.uur.LoginAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		http.Error(w, "Failed to Login or generate token", http.StatusInternalServerError)
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
	user, err := uh.ah.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	if err = uh.uur.LogoutUser(ctx, user.ID.String()); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
