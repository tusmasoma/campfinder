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
	"github.com/tusmasoma/campfinder/docker/back/usecase/mock"
)

func TestImageHandler_HandleImageGet(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockImageUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockImageUseCase, m1 *mock.MockAuthUseCase) {
				const layout = "2006-01-02T15:04:05Z"
				created, _ := time.Parse(layout, "0001-01-01T00:00:00Z")
				m.EXPECT().GetSpotImgURLBySpotID(
					gomock.Any(),
					"fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052",
				).Return(
					[]model.Image{
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
			iuc := mock.NewMockImageUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(iuc, auc)
			}

			handler := NewImageHandler(iuc, auc)
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
			m *mock.MockImageUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockImageUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}

				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().ImageCreate(
					gomock.Any(),
					uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					"https://hoge.com/hoge",
					user,
				).Return(nil)
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
			setup: func(m *mock.MockImageUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}

				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(&user, nil)
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
			iuc := mock.NewMockImageUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(iuc, auc)
			}

			handler := NewImageHandler(iuc, auc)
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
			m *mock.MockImageUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockImageUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}

				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(&user, nil)
				m.EXPECT().ImageDelete(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					user,
				).Return(nil)
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
			setup: func(m *mock.MockImageUseCase, m1 *mock.MockAuthUseCase) {
				superUser := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "super_user",
					Email:    "super_user@gmail.com",
					Password: "password123",
					IsAdmin:  true,
				}

				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(&superUser, nil)
				m.EXPECT().ImageDelete(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					superUser,
				).Return(nil)
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
			setup: func(m *mock.MockImageUseCase, m1 *mock.MockAuthUseCase) {
				user := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}

				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(&user, nil)
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
			setup: func(m *mock.MockImageUseCase, m1 *mock.MockAuthUseCase) {
				userWithoutAuth := model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				}

				m1.EXPECT().FetchUserFromContext(gomock.Any()).Return(&userWithoutAuth, nil)
				m.EXPECT().ImageDelete(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					userWithoutAuth,
				).Return(
					fmt.Errorf("Not authorized to update"),
				)
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
			iuc := mock.NewMockImageUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(iuc, auc)
			}

			handler := NewImageHandler(iuc, auc)
			recorder := httptest.NewRecorder()

			handler.HandleImageDelete(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
