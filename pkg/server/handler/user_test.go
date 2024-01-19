package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	cachemock "github.com/tusmasoma/campfinder/cache/mock"
	"github.com/tusmasoma/campfinder/db"
	dbmock "github.com/tusmasoma/campfinder/db/mock"
	authmock "github.com/tusmasoma/campfinder/internal/auth/mock"
	"github.com/tusmasoma/campfinder/pkg/server/handler/mock"
)

func TestUserHandler_HandleUserCreate(t *testing.T) {
	t.Error("テストが失敗しました。GitHub Actionsの特定のジョブが失敗した場合に、mainブランチへのマージを禁止する機能をテスト")
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserHandler,
			m1 *dbmock.MockUserRepository,
			m2 *cachemock.MockRedisRepository,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserHandler, m1 *dbmock.MockUserRepository, m2 *cachemock.MockRedisRepository) {
				t.Setenv("PRIVATE_KEY_PATH", "../../../.certificate/private_key.pem")
				m1.EXPECT().CheckIfUserExists(gomock.Any(), "test@gmail.com").Return(false, nil)
				m1.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				m2.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			in: func() *http.Request {
				userCreateReq := UserCreateRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			in: func() *http.Request {
				userCreateReq := UserCreateRequest{Email: "test@gmail.com"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: Username already exists",
			setup: func(
				m *mock.MockUserHandler,
				m1 *dbmock.MockUserRepository,
				m2 *cachemock.MockRedisRepository,
			) {
				t.Setenv("PRIVATE_KEY_PATH", "../../../.certificate/private_key.pem")
				m1.EXPECT().CheckIfUserExists(gomock.Any(), "test@gmail.com").Return(true, nil)
			},
			in: func() *http.Request {
				userCreateReq := UserCreateRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			repo := mock.NewMockUserHandler(ctrl)
			mockUserRepo := dbmock.NewMockUserRepository(ctrl)
			mockRedisRepo := cachemock.NewMockRedisRepository(ctrl)
			mockAuthHandler := authmock.NewMockHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockUserRepo, mockRedisRepo)
			}

			handler := NewUserHandler(mockUserRepo, mockRedisRepo, mockAuthHandler)
			recorder := httptest.NewRecorder()
			handler.HandleUserCreate(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				if token := recorder.Header().Get("Authorization"); token == "" || strings.TrimPrefix(token, "Bearer ") == "" {
					t.Fatalf("Expected Authorization header to be set")
				}
			}
		})
	}
}

func TestUserHandler_HandleUserLogin(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserHandler,
			m1 *dbmock.MockUserRepository,
			m2 *cachemock.MockRedisRepository,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserHandler, m1 *dbmock.MockUserRepository, m2 *cachemock.MockRedisRepository) {
				t.Setenv("PRIVATE_KEY_PATH", "../../../.certificate/private_key.pem")
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().GetUserByEmail(gomock.Any(), "test@gmail.com").Return(
					db.User{
						ID:       uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
					}, nil)
				m2.EXPECT().Exists(gomock.Any(), "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed").Return(false)
				m2.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			in: func() *http.Request {
				userLoginReq := UserLoginRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userLoginReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			in: func() *http.Request {
				userLoginReq := UserLoginRequest{Email: "test@gmail.com"}
				reqBody, _ := json.Marshal(userLoginReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: User alredy logined",
			setup: func(
				m *mock.MockUserHandler,
				m1 *dbmock.MockUserRepository,
				m2 *cachemock.MockRedisRepository,
			) {
				t.Setenv("PRIVATE_KEY_PATH", "../../../.certificate/private_key.pem")
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().GetUserByEmail(gomock.Any(), "test@gmail.com").Return(
					db.User{
						ID:       uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
					}, nil)
				m2.EXPECT().Exists(gomock.Any(), "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed").Return(true)
			},
			in: func() *http.Request {
				userLoginReq := UserLoginRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userLoginReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Fail: Passwords do not match",
			setup: func(
				m *mock.MockUserHandler,
				m1 *dbmock.MockUserRepository,
				m2 *cachemock.MockRedisRepository,
			) {
				t.Setenv("PRIVATE_KEY_PATH", "../../../.certificate/private_key.pem")
				passward, _ := db.PasswordEncrypt("password456")
				m1.EXPECT().GetUserByEmail(gomock.Any(), "test@gmail.com").Return(
					db.User{
						ID:       uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
					}, nil)
				m2.EXPECT().Exists(gomock.Any(), "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed").Return(false)
			},
			in: func() *http.Request {
				userLoginReq := UserLoginRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userLoginReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			repo := mock.NewMockUserHandler(ctrl)
			mockUserRepo := dbmock.NewMockUserRepository(ctrl)
			mockRedisRepo := cachemock.NewMockRedisRepository(ctrl)
			mockAuthHandler := authmock.NewMockHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockUserRepo, mockRedisRepo)
			}

			handler := NewUserHandler(mockUserRepo, mockRedisRepo, mockAuthHandler)
			recorder := httptest.NewRecorder()
			handler.HandleUserLogin(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				if token := recorder.Header().Get("Authorization"); token == "" || strings.TrimPrefix(token, "Bearer ") == "" {
					t.Fatalf("Expected Authorization header to be set")
				}
			}
		})
	}
}
