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

func TestCommentHandler_HandleCommentGet(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockCommentRepository,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *dbmock.MockCommentRepository,
			) {
				const layout = "2006-01-02T15:04:05Z"
				created, _ := time.Parse(layout, "0001-01-01T00:00:00Z")
				m.EXPECT().GetCommentBySpotID(gomock.Any(), "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052").Return(
					[]db.Comment{
						{
							ID:       uuid.New(),
							SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
							UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
							StarRate: 2,
							Text:     "いいスポットでした!!!",
							Created:  created,
						},
					}, nil,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/comment?spot_id=fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052", nil)
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
			repo := dbmock.NewMockCommentRepository(ctrl)
			mockAuthHandler := authmock.NewMockAuthHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			handler := NewCommentHandler(repo, mockAuthHandler)
			recorder := httptest.NewRecorder()

			handler.HandleCommentGet(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandler_HandleCommentCreate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockCommentRepository,
			m1 *authmock.MockAuthHandler,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *dbmock.MockCommentRepository,
				m1 *authmock.MockAuthHandler,
			) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
					IsAdmin:  false,
				}, nil)
				m.EXPECT().Create(gomock.Any(), db.Comment{
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 3,
					Text:     "いいスポットでした！",
				}).Return(nil)
			},
			in: func() *http.Request {
				commentCreateReq := CommentCreateRequest{
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					StarRate: 3,
					Text:     "いいスポットでした！",
				}
				reqBody, _ := json.Marshal(commentCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *dbmock.MockCommentRepository, m1 *authmock.MockAuthHandler) {
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
				commentCreateReq := CommentCreateRequest{
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					StarRate: 3,
				}
				reqBody, _ := json.Marshal(commentCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/create", bytes.NewBuffer(reqBody))
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
			repo := dbmock.NewMockCommentRepository(ctrl)
			mockAuthHandler := authmock.NewMockAuthHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockAuthHandler)
			}

			handler := NewCommentHandler(repo, mockAuthHandler)
			recorder := httptest.NewRecorder()

			handler.HandleCommentCreate(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandler_HandleCommentUpdate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockCommentRepository,
			m1 *authmock.MockAuthHandler,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *dbmock.MockCommentRepository,
				m1 *authmock.MockAuthHandler,
			) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
					IsAdmin:  false,
				}, nil)
				m.EXPECT().Update(gomock.Any(), db.Comment{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5,
					Text:     "いいスポットでした！!!",
				}).Return(nil)
			},
			in: func() *http.Request {
				commentUpdateReq := CommentUpdateRequest{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5,
					Text:     "いいスポットでした！!!",
				}
				reqBody, _ := json.Marshal(commentUpdateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "success: Super User",
			setup: func(
				m *dbmock.MockCommentRepository,
				m1 *authmock.MockAuthHandler,
			) {
				passward, _ := db.PasswordEncrypt("password123")
				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(db.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "super_user",
					Email:    "super_user@gmail.com",
					Password: passward,
					IsAdmin:  true,
				}, nil)
				m.EXPECT().Update(gomock.Any(), db.Comment{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5,
					Text:     "いいスポットでした！!!",
				}).Return(nil)
			},
			in: func() *http.Request {
				commentUpdateReq := CommentUpdateRequest{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5,
					Text:     "いいスポットでした！!!",
				}
				reqBody, _ := json.Marshal(commentUpdateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *dbmock.MockCommentRepository, m1 *authmock.MockAuthHandler) {
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
				commentUpdateReq := CommentUpdateRequest{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5,
				}
				reqBody, _ := json.Marshal(commentUpdateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: Not authorized to update",
			setup: func(
				m *dbmock.MockCommentRepository,
				m1 *authmock.MockAuthHandler,
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
				commentUpdateReq := CommentUpdateRequest{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5,
					Text:     "いいスポットでした！!!",
				}
				reqBody, _ := json.Marshal(commentUpdateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/update", bytes.NewBuffer(reqBody))
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
			repo := dbmock.NewMockCommentRepository(ctrl)
			mockAuthHandler := authmock.NewMockAuthHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockAuthHandler)
			}

			handler := NewCommentHandler(repo, mockAuthHandler)
			recorder := httptest.NewRecorder()

			handler.HandleCommentUpdate(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandler_HandleCommentDelete(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockCommentRepository,
			m1 *authmock.MockAuthHandler,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *dbmock.MockCommentRepository,
				m1 *authmock.MockAuthHandler,
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
				commentDeleteReq := CommentDeleteRequest{
					ID:     uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					UserID: uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
				}
				reqBody, _ := json.Marshal(commentDeleteReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/delete", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "success: Super User",
			setup: func(
				m *dbmock.MockCommentRepository,
				m1 *authmock.MockAuthHandler,
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
				commentDeleteReq := CommentDeleteRequest{
					ID:     uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					UserID: uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
				}
				reqBody, _ := json.Marshal(commentDeleteReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/delete", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *dbmock.MockCommentRepository, m1 *authmock.MockAuthHandler) {
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
				commentDeleteReq := CommentDeleteRequest{
					ID: uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
				}
				reqBody, _ := json.Marshal(commentDeleteReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/delete", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: Not authorized to update",
			setup: func(
				m *dbmock.MockCommentRepository,
				m1 *authmock.MockAuthHandler,
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
				commentDeleteReq := CommentDeleteRequest{
					ID:     uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					UserID: uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
				}
				reqBody, _ := json.Marshal(commentDeleteReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/delete", bytes.NewBuffer(reqBody))
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
			repo := dbmock.NewMockCommentRepository(ctrl)
			mockAuthHandler := authmock.NewMockAuthHandler(ctrl)

			if tt.setup != nil {
				tt.setup(repo, mockAuthHandler)
			}

			handler := NewCommentHandler(repo, mockAuthHandler)
			recorder := httptest.NewRecorder()

			handler.HandleCommentDelete(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v: %v", status, tt.wantStatus, tt.name)
			}
		})
	}
}
