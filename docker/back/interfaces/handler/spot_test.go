package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
	"github.com/tusmasoma/campfinder/docker/back/usecase/mock"
)

func TestSpotHandler_CreateSpot(t *testing.T) {
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
					&usecase.CreateSpotParams{
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
				).Return(nil)
			},
			in: func() *http.Request {
				spotCreateReq := CreateSpotRequest{
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
				spotCreateReq := CreateSpotRequest{
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
		{
			name: "Fail: create spot",
			setup: func(m *mock.MockSpotUseCase) {
				m.EXPECT().CreateSpot(
					gomock.Any(),
					&usecase.CreateSpotParams{
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
				).Return(fmt.Errorf("faile to create spot"))
			},
			in: func() *http.Request {
				spotCreateReq := CreateSpotRequest{
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
			wantStatus: http.StatusInternalServerError,
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

			handler.CreateSpot(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestSpotHandler_BatchCreateSpots(t *testing.T) {
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
				m.EXPECT().BatchCreateSpots(
					gomock.Any(),
					&usecase.BatchCreateSpotParams{
						Spots: []usecase.CreateSpotParams{
							{
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
					},
				).Return(nil)
			},
			in: func() *http.Request {
				spotBatchCreateReq := BatchCreateSpotsRequest{
					Spots: []CreateSpotRequest{
						{
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
				}
				reqBody, _ := json.Marshal(spotBatchCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/spot/batchcreate", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
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

			handler.BatchCreateSpots(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestSpotHandler_ListSpots(t *testing.T) {
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
				m.EXPECT().ListSpots(gomock.Any(), []string{"campsite", "spa"}).Return(
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

			handler.ListSpots(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestSpotHandler_GetSpot(t *testing.T) {
	t.Parallel()

	spot := model.Spot{
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
	}
	jsonData, _ := json.MarshalIndent(spot, "", "    ")

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotUseCase,
		)
		spotID     string
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			setup: func(m *mock.MockSpotUseCase) {
				m.EXPECT().GetSpot(
					gomock.Any(),
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
				).Return(spot)
			},
			spotID:     "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			wantStatus: http.StatusOK,
			wantBody:   string(jsonData),
		},
		{
			name: "success: getspot method return nil",
			setup: func(m *mock.MockSpotUseCase) {
				m.EXPECT().GetSpot(
					gomock.Any(),
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
				).Return(model.Spot{})
			},
			spotID:     "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			wantStatus: http.StatusOK,
			wantBody:   `{"spot":{}}`,
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

			r := chi.NewRouter()
			r.Get("/api/spot/{spotID}", handler.GetSpot)

			req, _ := http.NewRequest(http.MethodGet, "/api/spot/"+tt.spotID, nil)
			r.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
