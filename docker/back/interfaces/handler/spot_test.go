package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/usecase/mock"
)

func TestSpotHandler_HandleSpotCreate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *mock.MockSpotUseCase,
			) {
				m.EXPECT().CreateSpot(
					gomock.Any(),
					"campsite",
					"旭川市21世紀の森ふれあい広場",
					"北海道旭川市東旭川町瑞穂4288",
					43.7172721,
					142.6674615,
					"2022年5月1日(日)〜11月30日(水)",
					"0166-76-2108",
					"有料。ログハウス大人290円〜750円、高校生以下180〜460円",
					"旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。",
					"/static/img/campsiteflag.jpeg",
				).Return(nil)
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
			repo := mock.NewMockSpotUseCase(ctrl)

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
			m *mock.MockSpotUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(
				m *mock.MockSpotUseCase,
			) {
				m.EXPECT().ListSpots(gomock.Any(), []string{"campsite"}, "0c0000e0-c00f-0dac-00ef-d00ab0ea0fed").Return(
					[]model.Spot{
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
						{
							ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab8ea5fde"),
							Category:    "campsite",
							Name:        "とままえ夕陽ヶ丘未来港公園",
							Address:     "北海道苫前郡苫前町字栄浜313",
							Lat:         44.3153234,
							Lng:         141.6563455,
							Period:      "管理棟は7月中旬～8月中旬",
							Phone:       "0164-64-2212",
							Price:       "不明",
							Description: "とままえ夕陽ヶ丘公園は、日本海に面した位置にある開放感あふれる公園です。",
							IconPath:    "/static/img/campsiteflag.jpeg",
						},
						{
							ID:          uuid.MustParse("0c0000e0-c00f-0dac-00ef-d00ab0ea0fed"),
							Category:    "spa",
							Name:        "奥の湯",
							Address:     "北海道川上郡弟子屈町字屈斜路",
							Lat:         43.566446,
							Lng:         144.3091296,
							Period:      "24時間",
							Phone:       "-",
							Price:       "無料",
							Description: "奥の湯は、札幌市北34条駅から徒歩0分という便利なロケーションにある銭湯です。",
							IconPath:    "/static/img/spaflag.jpeg",
						},
					},
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodGet,
					"/api/spot?category=campsite&spot_id=0c0000e0-c00f-0dac-00ef-d00ab0ea0fed",
					nil,
				)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "success: only category",
			setup: func(
				m *mock.MockSpotUseCase,
			) {
				m.EXPECT().ListSpots(gomock.Any(), []string{"campsite", "spa"}, "").Return(
					[]model.Spot{
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
						{
							ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab8ea5fde"),
							Category:    "campsite",
							Name:        "とままえ夕陽ヶ丘未来港公園",
							Address:     "北海道苫前郡苫前町字栄浜313",
							Lat:         44.3153234,
							Lng:         141.6563455,
							Period:      "管理棟は7月中旬～8月中旬",
							Phone:       "0164-64-2212",
							Price:       "不明",
							Description: "とままえ夕陽ヶ丘公園は、日本海に面した位置にある開放感あふれる公園です。",
							IconPath:    "/static/img/campsiteflag.jpeg",
						},
						{
							ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab8ea5fdf"),
							Category:    "spa",
							Name:        "奥の湯",
							Address:     "北海道川上郡弟子屈町字屈斜路",
							Lat:         43.566446,
							Lng:         144.3091296,
							Period:      "24時間",
							Phone:       "-",
							Price:       "無料",
							Description: "奥の湯は、札幌市北34条駅から徒歩0分という便利なロケーションにある銭湯です。",
							IconPath:    "/static/img/spaflag.jpeg",
						},
					},
				)
			},

			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodGet,
					"/api/spot?category=campsite&category=spa",
					nil,
				)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "success: only spot_id",
			setup: func(
				m *mock.MockSpotUseCase,
			) {
				m.EXPECT().ListSpots(gomock.Any(), nil, "0c0000e0-c00f-0dac-00ef-d00ab0ea0fed").Return(
					[]model.Spot{
						{
							ID:          uuid.MustParse("0c0000e0-c00f-0dac-00ef-d00ab0ea0fed"),
							Category:    "spa",
							Name:        "奥の湯",
							Address:     "北海道川上郡弟子屈町字屈斜路",
							Lat:         43.566446,
							Lng:         144.3091296,
							Period:      "24時間",
							Phone:       "-",
							Price:       "無料",
							Description: "奥の湯は、札幌市北34条駅から徒歩0分という便利なロケーションにある銭湯です。",
							IconPath:    "/static/img/spaflag.jpeg",
						},
					},
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(
					http.MethodGet,
					"/api/spot?spot_id=0c0000e0-c00f-0dac-00ef-d00ab0ea0fed",
					nil,
				)
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
			repo := mock.NewMockSpotUseCase(ctrl)

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
