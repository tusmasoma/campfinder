package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/db"
	dbmock "github.com/tusmasoma/campfinder/db/mock"
	authmock "github.com/tusmasoma/campfinder/internal/auth/mock"
)

func TestImageHandler_HandleImageGet(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockImageRepository,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *dbmock.MockImageRepository,
			) {
				const layout = "2006-01-02T15:04:05Z"
				created, _ := time.Parse(layout, "0001-01-01T00:00:00Z")
				m.EXPECT().GetSpotImgURLBySpotID(gomock.Any(), "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052").Return(
					[]db.Image{
						{
							ID:      uuid.New(),
							SpotID:  uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
							UserID:  uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
							URL:     "https://hoge.com/hoge",
							Created: created,
						},
					}, nil,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/img?spot_id=fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052", nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := dbmock.NewMockImageRepository(ctrl)
			mockAuthHandler := authmock.NewMockHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			handler := NewImageHandler(repo, mockAuthHandler)
			recorder := httptest.NewRecorder()

			handler.HandleImageGet(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestImageHandler_HandleImageCreate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockImageRepository,
			m1 *authmock.MockHandler,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *dbmock.MockImageRepository,
				m1 *authmock.MockHandler,
			) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
					IsAdmin:  false,
				}, nil)
				m.EXPECT().Create(gomock.Any(), db.Image{
					SpotID: uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID: uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					URL:    "https://hoge.com/hoge",
				}).Return(nil)
			},
			in: func() *http.Request {
				imgCreateReq := ImageCreateRequest{
					SpotID: uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					URL:    "https://hoge.com/hoge",
				}
				reqBody, _ := json.Marshal(imgCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/img/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *dbmock.MockImageRepository, m1 *authmock.MockHandler) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
					IsAdmin:  false,
				}, nil)
			},
			in: func() *http.Request {
				imgCreateReq := ImageCreateRequest{
					SpotID: uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
				}
				reqBody, _ := json.Marshal(imgCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/img/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := dbmock.NewMockImageRepository(ctrl)
			mockAuthHandler := authmock.NewMockHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockAuthHandler)
			}

			handler := NewImageHandler(repo, mockAuthHandler)
			recorder := httptest.NewRecorder()

			handler.HandleImageCreate(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestImageHandler_HandleImageDelete(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockImageRepository,
			m1 *authmock.MockHandler,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *dbmock.MockImageRepository,
				m1 *authmock.MockHandler,
			) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
					IsAdmin:  false,
				}, nil)
				m.EXPECT().Delete(gomock.Any(), "31894386-3e60-45a8-bc67-f46b72b42554").Return(nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/img/delete?id=31894386-3e60-45a8-bc67-f46b72b42554&user_id=f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "success: Super User",
			setup: func(
				m *dbmock.MockImageRepository,
				m1 *authmock.MockHandler,
			) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "super_user",
					Email:    "super_user@gmail.com",
					Password: passward,
					IsAdmin:  true,
				}, nil)
				m.EXPECT().Delete(gomock.Any(), "31894386-3e60-45a8-bc67-f46b72b42554").Return(nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/img/delete?id=31894386-3e60-45a8-bc67-f46b72b42554&user_id=f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *dbmock.MockImageRepository, m1 *authmock.MockHandler) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
					IsAdmin:  false,
				}, nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/img/delete?id=31894386-3e60-45a8-bc67-f46b72b42554",
					nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: Not authorized to update",
			setup: func(
				m *dbmock.MockImageRepository,
				m1 *authmock.MockHandler,
			) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
					IsAdmin:  false,
				}, nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/img/delete?id=31894386-3e60-45a8-bc67-f46b72b42554&user_id=f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					nil)
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
			repo := dbmock.NewMockImageRepository(ctrl)
			mockAuthHandler := authmock.NewMockHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockAuthHandler)
			}

			handler := NewImageHandler(repo, mockAuthHandler)
			recorder := httptest.NewRecorder()

			handler.HandleImageDelete(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v: %v", status, tt.wantStatus, tt.name)
			}
		})
	}
}
