package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	cachemock "github.com/tusmasoma/campfinder/cache/mock"
	dbmock "github.com/tusmasoma/campfinder/db/mock"
	"github.com/tusmasoma/campfinder/pkg/server/handler/mock"
)

func TestUserHandler_HandleUserCreate(t *testing.T) {

	t.Parallel()
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
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := mock.NewMockUserHandler(ctrl)
			mockUserRepo := dbmock.NewMockUserRepository(ctrl)
			mockRedisRepo := cachemock.NewMockRedisRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockUserRepo, mockRedisRepo)
			}

			handler := NewUserHandler(mockUserRepo, mockRedisRepo)
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
