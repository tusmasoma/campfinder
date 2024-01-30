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
