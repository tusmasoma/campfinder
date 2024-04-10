package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository/mock"
)

type SpotCreateArg struct {
	ctx         context.Context
	category    string
	name        string
	address     string
	lat         float64
	lng         float64
	period      string
	phone       string
	price       string
	description string
	iconPath    string
}

type ListSpotsArg struct {
	ctx        context.Context
	categories []string
}

type GetSpotArg struct {
	ctx    context.Context
	spotID string
}

func TestSpotUseCase_SpotCreate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotRepository,
		)
		arg     SpotCreateArg
		wantErr error
	}{
		{
			name: "sccess",
			setup: func(m *mock.MockSpotRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{
						{Field: "Lat", Value: 43.7172721},
						{Field: "Lng", Value: 142.6674615},
					},
				).Return([]model.Spot{}, nil)
				m.EXPECT().Create(
					gomock.Any(),
					model.Spot{
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
			arg: SpotCreateArg{
				ctx:         context.Background(), // コンテキストを適切に設定
				category:    "campsite",
				name:        "旭川市21世紀の森ふれあい広場",
				address:     "北海道旭川市東旭川町瑞穂4288",
				lat:         43.7172721,
				lng:         142.6674615,
				period:      "2022年5月1日(日)〜11月30日(水)",
				phone:       "0166-76-2108",
				price:       "有料。ログハウス大人290円〜750円、高校生以下180〜460円",
				description: "旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。",
				iconPath:    "/static/img/campsiteflag.jpeg",
			},
			wantErr: nil,
		},
		{
			name: "fail: already exists",
			setup: func(m *mock.MockSpotRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{
						{Field: "Lat", Value: 43.7172721},
						{Field: "Lng", Value: 142.6674615},
					},
				).Return([]model.Spot{
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
				}, nil)
			},
			arg: SpotCreateArg{
				ctx:         context.Background(), // コンテキストを適切に設定
				category:    "campsite",
				name:        "旭川市21世紀の森ふれあい広場",
				address:     "北海道旭川市東旭川町瑞穂4288",
				lat:         43.7172721,
				lng:         142.6674615,
				period:      "2022年5月1日(日)〜11月30日(水)",
				phone:       "0166-76-2108",
				price:       "有料。ログハウス大人290円〜750円、高校生以下180〜460円",
				description: "旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。",
				iconPath:    "/static/img/campsiteflag.jpeg",
			},
			wantErr: fmt.Errorf("already exists"),
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			sr := mock.NewMockSpotRepository(ctrl)
			cr := mock.NewMockCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(sr)
			}

			usecase := NewSpotUseCase(sr, cr)

			err := usecase.CreateSpot(
				tt.arg.ctx,
				tt.arg.category,
				tt.arg.name,
				tt.arg.address,
				tt.arg.lat,
				tt.arg.lng,
				tt.arg.period,
				tt.arg.phone,
				tt.arg.price,
				tt.arg.description,
				tt.arg.iconPath,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("SpotCreate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("SpotCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSpotUseCase_ListSpots(t *testing.T) {
	t.Parallel()
	campsite := model.Spot{
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
	spa := model.Spot{
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
	}
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotRepository,
			m1 *mock.MockCacheRepository,
		)
		arg  ListSpotsArg
		want []model.Spot
	}{
		{
			name: "success",
			setup: func(m *mock.MockSpotRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Category", Value: "campsite"}},
				).Return([]model.Spot{campsite}, nil)
				m1.EXPECT().Set(
					gomock.Any(),
					"spots_campsite",
					[]model.Spot{campsite},
				).Return(nil)
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Category", Value: "spa"}},
				).Return([]model.Spot{spa}, nil)
				m1.EXPECT().Set(
					gomock.Any(),
					"spots_spa",
					[]model.Spot{spa},
				).Return(nil)
			},
			arg: ListSpotsArg{
				ctx:        context.Background(),
				categories: []string{"campsite", "spa"},
			},
			want: []model.Spot{campsite, spa},
		},
		{
			name: "fail: get spot from db",
			setup: func(m *mock.MockSpotRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Category", Value: "campsite"}},
				).Return([]model.Spot{}, fmt.Errorf("fail to get spot from db"))
				m1.EXPECT().Get(gomock.Any(), "spots_campsite").Return([]model.Spot{campsite}, nil)
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Category", Value: "spa"}},
				).Return([]model.Spot{}, fmt.Errorf("fail to get spot from db"))
				m1.EXPECT().Get(gomock.Any(), "spots_spa").Return([]model.Spot{spa}, nil)
			},
			arg: ListSpotsArg{
				ctx:        context.Background(),
				categories: []string{"campsite", "spa"},
			},
			want: []model.Spot{campsite, spa},
		},
		{
			name: "fail: set master data",
			setup: func(m *mock.MockSpotRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Category", Value: "campsite"}},
				).Return([]model.Spot{campsite}, nil)
				m1.EXPECT().Set(
					gomock.Any(),
					"spots_campsite",
					[]model.Spot{campsite},
				).Return(fmt.Errorf("fail to set in cache"))
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Category", Value: "spa"}},
				).Return([]model.Spot{spa}, nil)
				m1.EXPECT().Set(
					gomock.Any(),
					"spots_spa",
					[]model.Spot{spa},
				).Return(fmt.Errorf("fail to set in cache"))
			},
			arg: ListSpotsArg{
				ctx:        context.Background(),
				categories: []string{"campsite", "spa"},
			},
			want: []model.Spot{campsite, spa},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			sr := mock.NewMockSpotRepository(ctrl)
			cr := mock.NewMockCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(sr, cr)
			}

			usecase := NewSpotUseCase(sr, cr)

			spots := usecase.ListSpots(tt.arg.ctx, tt.arg.categories)

			if !reflect.DeepEqual(spots, tt.want) {
				t.Errorf("ListSpots() \n got = %v,\n want %v", spots, tt.want)
			}
		})
	}
}

func TestSpotUseCase_GetSpot(t *testing.T) {
	t.Parallel()
	campsite := model.Spot{
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
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotRepository,
			m1 *mock.MockCacheRepository,
		)
		arg  GetSpotArg
		want model.Spot
	}{
		{
			name: "success",
			setup: func(m *mock.MockSpotRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().Get(gomock.Any(), "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed").Return(&campsite, nil)
			},
			arg: GetSpotArg{
				ctx:    context.Background(),
				spotID: "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			},
			want: campsite,
		},
		{
			name: "fail: fail to get spot form db. but, success to get spot from cache.",
			setup: func(m *mock.MockSpotRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().Get(gomock.Any(), "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed").Return(nil, fmt.Errorf("fail to get spot from db"))
				m1.EXPECT().Scan(gomock.Any(), "spots_*").Return([]string{"spots_campsite", "spots_spa"}, nil)
				m1.EXPECT().Get(gomock.Any(), "spots_campsite").Return([]model.Spot{campsite}, nil)
				m1.EXPECT().Get(gomock.Any(), "spots_spa").Return([]model.Spot{}, nil)
			},
			arg: GetSpotArg{
				ctx:    context.Background(),
				spotID: "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			},
			want: campsite,
		},
		{
			name: "fail: fail to get spot form db. and, does not exists the spot from cache.",
			setup: func(m *mock.MockSpotRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().Get(gomock.Any(), "5c5323e9-c78f-4dac-94ef-d34ab5ea8def").Return(nil, fmt.Errorf("fail to get spot from db"))
				m1.EXPECT().Scan(gomock.Any(), "spots_*").Return([]string{"spots_campsite", "spots_spa"}, nil)
				m1.EXPECT().Get(gomock.Any(), "spots_campsite").Return([]model.Spot{campsite}, nil)
				m1.EXPECT().Get(gomock.Any(), "spots_spa").Return([]model.Spot{}, nil)
			},
			arg: GetSpotArg{
				ctx:    context.Background(),
				spotID: "5c5323e9-c78f-4dac-94ef-d34ab5ea8def",
			},
			want: model.Spot{},
		},
		{
			name: "fail: fail to get spot form db. and, scan",
			setup: func(m *mock.MockSpotRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().Get(gomock.Any(), "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed").Return(nil, fmt.Errorf("fail to get spot from db"))
				m1.EXPECT().Scan(gomock.Any(), "spots_*").Return([]string{"spots_campsite", "spots_spa"}, fmt.Errorf("fail to scan"))
			},
			arg: GetSpotArg{
				ctx:    context.Background(),
				spotID: "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			},
			want: model.Spot{},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			sr := mock.NewMockSpotRepository(ctrl)
			cr := mock.NewMockCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(sr, cr)
			}

			usecase := NewSpotUseCase(sr, cr)

			spots := usecase.GetSpot(tt.arg.ctx, tt.arg.spotID)

			if !reflect.DeepEqual(spots, tt.want) {
				t.Errorf("GetSpot() \n got = %v,\n want %v", spots, tt.want)
			}
		})
	}
}
