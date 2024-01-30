package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository"
	"github.com/tusmasoma/campfinder/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUserAndGenerateToken(ctx context.Context, email string, passward string) (string, error)
}

type userUseCase struct {
	ur repository.UserRepository
}

func NewUserUseCase(ur repository.UserRepository) UserUseCase {
	return &userUseCase{
		ur: ur,
	}
}

func (uuc *userUseCase) CreateUserAndGenerateToken(ctx context.Context, email string, passward string) (string, error) {
	// ユーザー登録
	user, err := uuc.CreateUser(ctx, email, passward)
	if err != nil {
		log.Printf("Internal server error while creating user")
		return "", err
	}

	jwt, _ := auth.GenerateToken(user.ID.String(), user.Email)
	return jwt, nil
}

func (uuc *userUseCase) CreateUser(ctx context.Context, email string, passward string) (*model.User, error) {
	exists, err := uuc.ur.CheckIfUserExists(ctx, email)
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return nil, err
	}
	if exists {
		log.Printf("User with this name already exists - status: %d", http.StatusConflict)
		return nil, fmt.Errorf("user with this name already exists")
	}

	var user model.User
	user.Email = email
	user.Name = ExtractUsernameFromEmail(email)
	password, err := PasswordEncrypt(passward)
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return nil, err
	}
	user.Password = password

	if err = uuc.ur.Create(ctx, &user); err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}
	return &user, nil
}

func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ExtractUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
