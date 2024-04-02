//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tusmasoma/campfinder/docker/back/usecase"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	uur usecase.UserUseCase
	auc usecase.AuthUseCase
}

func NewUserHandler(uur usecase.UserUseCase, auc usecase.AuthUseCase) UserHandler {
	return &userHandler{
		uur: uur,
		auc: auc,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody CreateUserRequest
	if ok := isValidCreateUserRequest(r.Body, &requestBody); !ok {
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

func isValidCreateUserRequest(body io.ReadCloser, requestBody *CreateUserRequest) bool {
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

func (uh *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody LoginRequest
	if ok := isValidLoginRequest(r.Body, &requestBody); !ok {
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

func isValidLoginRequest(body io.ReadCloser, requestBody *LoginRequest) bool {
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

func (uh *userHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := uh.auc.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get UserInfo from context: %v", err), http.StatusInternalServerError)
		return
	}

	if err = uh.uur.LogoutUser(ctx, user.ID.String()); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
