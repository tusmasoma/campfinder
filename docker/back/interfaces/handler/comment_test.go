package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
	"github.com/tusmasoma/campfinder/docker/back/usecase/mock"
)

func TestCommentHandler_ListComments(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				const layout = "2006-01-02T15:04:05Z"
				created, _ := time.Parse(layout, "0001-01-01T00:00:00Z")
				m.EXPECT().ListComments(
					gomock.Any(),
					"fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052",
				).Return(
					[]model.Comment{
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
			cuc := mock.NewMockCommentUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(cuc, auc)
			}

			handler := NewCommentHandler(cuc, auc)
			recorder := httptest.NewRecorder()

			handler.ListComments(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandler_CreateComment(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().CreateComment(
					gomock.Any(),
					&usecase.CreateCommentParams{
						UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
						SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
						StarRate: 3.0,
						Text:     "いいスポットでした！",
					},
				).Return(nil)
			},
			in: func() *http.Request {
				commentCreateReq := CreateCommentRequest{
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					StarRate: 3.0,
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
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
			},
			in: func() *http.Request {
				commentCreateReq := CreateCommentRequest{
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					StarRate: 3.0,
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
			cuc := mock.NewMockCommentUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(cuc, auc)
			}

			handler := NewCommentHandler(cuc, auc)
			recorder := httptest.NewRecorder()

			handler.CreateComment(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandler_BatchCreateComments(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "",
					Password: "",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().BatchCreateComments(
					gomock.Any(),
					&usecase.BatchCreateCommentsParams{
						Comments: []usecase.CreateCommentParams{
							{
								UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
								SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
								StarRate: 5.0,
								Text:     "最高なスポットでした！",
							},
							{
								UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
								SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
								StarRate: 3.0,
								Text:     "いいスポットでした！",
							},
						},
					},
				).Return(nil)
			},
			in: func() *http.Request {
				commentCreateReq := BatchCreateCommentsRequest{
					Comments: []CreateCommentRequest{
						{
							SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
							StarRate: 5.0,
							Text:     "最高なスポットでした！",
						},
						{
							SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
							StarRate: 3.0,
							Text:     "いいスポットでした！",
						},
					},
				}
				reqBody, _ := json.Marshal(commentCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/batch_create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "",
					Password: "",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
			},
			in: func() *http.Request {
				commentCreateReq := BatchCreateCommentsRequest{
					Comments: []CreateCommentRequest{
						{
							SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
							StarRate: 3.0,
						},
					},
				}
				reqBody, _ := json.Marshal(commentCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/comment/batch_create", bytes.NewBuffer(reqBody))
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
			cuc := mock.NewMockCommentUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(cuc, auc)
			}

			handler := NewCommentHandler(cuc, auc)
			recorder := httptest.NewRecorder()

			handler.BatchCreateComments(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandler_UpdateComment(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().UpdateComment(
					gomock.Any(),
					uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					5.0,
					"いいスポットでした！!!",
					user,
				).Return(nil)
			},
			in: func() *http.Request {
				commentUpdateReq := UpdateCommentRequest{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5.0,
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
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				superUser := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "super_user",
					Email:    "super_user@gmail.com",
					Password: "password123",
					IsAdmin:  true,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&superUser, nil)
				m.EXPECT().UpdateComment(
					gomock.Any(),
					uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					5.0,
					"いいスポットでした！!!",
					superUser,
				).Return(nil)
			},
			in: func() *http.Request {
				commentUpdateReq := UpdateCommentRequest{
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
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
			},
			in: func() *http.Request {
				commentUpdateReq := UpdateCommentRequest{
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
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "super_user",
					Email:    "super_user@gmail.com",
					Password: "password123",
					IsAdmin:  true,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().UpdateComment(
					gomock.Any(),
					uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					5.0,
					"いいスポットでした！!!",
					user,
				).Return(
					fmt.Errorf("Not authorized to update"),
				)
			},
			in: func() *http.Request {
				commentUpdateReq := UpdateCommentRequest{
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
			cuc := mock.NewMockCommentUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(cuc, auc)
			}

			handler := NewCommentHandler(cuc, auc)
			recorder := httptest.NewRecorder()

			handler.UpdateComment(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandler_DeleteComment(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().DeleteComment(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					user,
				).Return(nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/comment/delete?id=31894386-3e60-45a8-bc67-f46b72b42554&user_id=f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "success: Super User",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				superUser := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "super_user",
					Email:    "super_user@gmail.com",
					Password: "password123",
					IsAdmin:  true,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&superUser, nil)
				m.EXPECT().DeleteComment(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					superUser,
				).Return(nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/comment/delete?id=31894386-3e60-45a8-bc67-f46b72b42554&user_id=f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/comment/delete?id=31894386-3e60-45a8-bc67-f46b72b42554",
					nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: Not authorized to update",
			setup: func(m *mock.MockCommentUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().DeleteComment(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
					"f6db2530-cd9b-4ac1-8dc1-38c795e61234",
					user,
				).Return(
					fmt.Errorf("Not authorized to delete"),
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodDelete,
					"/api/comment/delete?id=31894386-3e60-45a8-bc67-f46b72b42554&user_id=f6db2530-cd9b-4ac1-8dc1-38c795e61234",
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
			cuc := mock.NewMockCommentUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(cuc, auc)
			}

			handler := NewCommentHandler(cuc, auc)
			recorder := httptest.NewRecorder()

			handler.DeleteComment(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
