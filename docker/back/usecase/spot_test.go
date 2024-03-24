package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
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

type SpotGetArg struct {
	ctx        context.Context
	categories []string
	spotID     string
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
				m.EXPECT().CheckIfSpotExists(gomock.Any(), 43.7172721, 142.6674615).Return(false, nil)
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
				m.EXPECT().CheckIfSpotExists(gomock.Any(), 43.7172721, 142.6674615).Return(true, nil)
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
			wantErr: fmt.Errorf("user with this name already exists"),
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

			usecase := NewSpotUseCase(repo)

			err := usecase.SpotCreate(
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

func TestSpotUseCase_SpotGet(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockSpotRepository,
		)
		arg  SpotGetArg
		want []model.Spot
	}{
		{
			name: "success",
			setup: func(m *mock.MockSpotRepository) {
				m.EXPECT().GetSpotByCategory(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(
					// 無名関数の引数と戻り値はモックメソッドに揃える
					func(_ context.Context, category string, _ ...interface{}) ([]model.Spot, error) {
						switch category {
						case "campsite":
							return []model.Spot{
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
							}, nil
						case "spa":
							return []model.Spot{
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
							}, nil
						}
						return nil, nil
					},
				).AnyTimes()
				m.EXPECT().GetSpotByID(
					gomock.Any(),
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8def",
				).Return(
					model.Spot{
						ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8def"),
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
					nil,
				)
			},
			arg: SpotGetArg{
				ctx:        context.Background(),
				categories: []string{"campsite", "spa"},
				spotID:     "5c5323e9-c78f-4dac-94ef-d34ab5ea8def",
			},
			want: []model.Spot{
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
				{
					ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8def"),
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
			},
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

			usecase := NewSpotUseCase(repo)

			spots := usecase.SpotGet(tt.arg.ctx, tt.arg.categories, tt.arg.spotID)

			if !reflect.DeepEqual(spots, tt.want) {
				t.Errorf("SpotGet() \n got = %v,\n want %v", spots, tt.want)
			}
		})
	}
}
