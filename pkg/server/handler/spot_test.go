//lint:ignore testpackage
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/db"
	"github.com/tusmasoma/campfinder/db/mock"
)

func TestSpotHandler_HandleSpotCreate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotRepository,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *mock.MockSpotRepository,
			) {
				m.EXPECT().CheckIfSpotExists(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			in: func() *http.Request {
				spotCreateReq := SpotCreateRequest{
					Category:    "campsite",
					Name:        "旭川市21世紀の森ふれあい広場",
					Address:     "北海道旭川市東旭川町瑞穂4288",
					Lat:         43.7172721,
					Lng:         142.6674615,
					Period:      "2022年5月1日(日)〜11月30日(水)",
					Phone:       "0166-76-2108",
					Price:       "有料。ログハウス大人290円〜750円、高校生以下180〜460円",
					Description: "旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。",
					IconPath:    "/static/img/campsiteflag.jpeg",
				}
				reqBody, _ := json.Marshal(spotCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/spot/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			in: func() *http.Request {
				spotCreateReq := SpotCreateRequest{
					Category: "campsite",
					Name:     "旭川市21世紀の森ふれあい広場",
					Address:  "北海道旭川市東旭川町瑞穂4288",
				}
				reqBody, _ := json.Marshal(spotCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/spot/create", bytes.NewBuffer(reqBody))
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
			repo := mock.NewMockSpotRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			handler := NewSpotHandler(repo)
			recorder := httptest.NewRecorder()

			handler.HandleSpotCreate(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestSpotHandler_HandleSpotGet(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotRepository,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *mock.MockSpotRepository,
			) {
				m.EXPECT().GetSpotByCategory(gomock.Any(), "campsite").Return(
					[]db.Spot{
						{
							ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
							Category:    "campsite",
							Name:        "旭川市21世紀の森ふれあい広場",
							Address:     "北海道旭川市東旭川町瑞穂4288",
							Lat:         43.7172721,
							Lng:         142.6674615,
							Period:      "2022年5月1日(日)〜11月30日(水)",
							Phone:       "0166-76-2108",
							Price:       "有料。ログハウス大人290円〜750円、高校生以下180〜460円",
							Description: "旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。",
							IconPath:    "/static/img/campsiteflag.jpeg",
						},
					}, nil,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/spot?category=campsite", nil)
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
			repo := mock.NewMockSpotRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			handler := NewSpotHandler(repo)
			recorder := httptest.NewRecorder()

			handler.HandleSpotGet(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
