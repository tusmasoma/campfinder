package infra

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/domain/model"
	"gotest.tools/assert"
)

type CheckIfSpotExistsArg struct {
	ctx context.Context
	lat float64
	lng float64
}

type GetSpotByIDArg struct {
	ctx context.Context
	id  string
}

type GetSpotByCategoryArg struct {
	ctx      context.Context
	category string
}

type SpotCreateArg struct {
	ctx  context.Context
	spot model.Spot
}

type SpotUpdateArg struct {
	ctx  context.Context
	spot model.Spot
}

type SpotDeleteArg struct {
	ctx  context.Context
	spot model.Spot
}

func TestSpotRepo_CheckIfSpotExists(t *testing.T) {
	patterns := []struct {
		name       string
		setup      func(db *sql.DB)
		in         CheckIfSpotExistsArg
		wantExists bool
	}{
		{
			name: "Success",
			in: CheckIfSpotExistsArg{
				ctx: context.Background(),
				lat: 43.7172721,
				lng: 142.6674615,
			},
			wantExists: true,
		},
		{
			name: "Fail",
			in: CheckIfSpotExistsArg{
				ctx: context.Background(),
				lat: 00.0000000,
				lng: 000.0000000,
			},
			wantExists: false,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewSpotRepository(db)

			exists, _ := repo.CheckIfSpotExists(tt.in.ctx, tt.in.lat, tt.in.lng)

			assert.Equal(t, tt.wantExists, exists)
		})
	}
}

func TestSpotRepo_GetSpotByID(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    GetSpotByIDArg
		want  struct {
			spot model.Spot
			err  error
		}
	}{
		{
			name: "Success",
			in: GetSpotByIDArg{
				ctx: context.Background(),
				id:  "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			},
			want: struct {
				spot model.Spot
				err  error
			}{
				spot: model.Spot{
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
				err: nil,
			},
		},
		{
			name: "Fail",
			in: GetSpotByIDArg{
				ctx: context.Background(),
				id:  "00000000-0000-0000-0000-000000000000",
			},
			want: struct {
				spot model.Spot
				err  error
			}{
				spot: model.Spot{},
				err:  sql.ErrNoRows,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewSpotRepository(db)

			spot, err := repo.GetSpotByID(tt.in.ctx, tt.in.id)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetSpotByID() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetSpotByID() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(spot, tt.want.spot) {
				t.Errorf("GetSpotByID() \n got = %v,\n want %v", spot, tt.want.spot)
			}
		})
	}
}

func TestSpotRepo_GetSpotByCategory(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    GetSpotByCategoryArg
		want  struct {
			spots []model.Spot
			err   error
		}
	}{
		{
			name: "Success",
			in: GetSpotByCategoryArg{
				ctx:      context.Background(),
				category: "campsite",
			},
			want: struct {
				spots []model.Spot
				err   error
			}{
				spots: []model.Spot{
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
				},
				err: nil,
			},
		},
		{
			name: "Fail",
			in: GetSpotByCategoryArg{
				ctx:      context.Background(),
				category: "spa",
			},
			want: struct {
				spots []model.Spot
				err   error
			}{
				spots: nil,
				err:   nil,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewSpotRepository(db)

			spots, err := repo.GetSpotByCategory(tt.in.ctx, tt.in.category)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetSpotByCategory() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetSpotByCategory() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(spots, tt.want.spots) {
				t.Errorf("GetSpotByCategory() \n got = %v,\n want %v", spots, tt.want.spots)
			}
		})
	}
}

func TestSpotRepo_Create(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name    string
		setup   func(db *sql.DB)
		in      SpotCreateArg
		wantErr error
	}{
		{
			name: "Success",
			in: SpotCreateArg{
				ctx: context.Background(),
				spot: model.Spot{
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
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewSpotRepository(db)

			err := repo.Create(tt.in.ctx, tt.in.spot)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				var exists bool
				exists, err = repo.CheckIfSpotExists(tt.in.ctx, tt.in.spot.Lat, tt.in.spot.Lng)
				if !exists {
					t.Errorf("Create() is successful but there is not spot in db:  err = %v", err)
				}
			}
		})
	}
}

func TestSpotRepo_Update(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    SpotUpdateArg
		want  struct {
			Period string
			err    error
		}
	}{
		{
			name: "Success",
			in: SpotUpdateArg{
				ctx: context.Background(),
				spot: model.Spot{
					ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8abc"),
					Category:    "campsite",
					Name:        "千代田の丘キャンプ場",
					Address:     "北海道上川郡美瑛町字水沢春日台第一",
					Lat:         43.5436008,
					Lng:         142.4912747,
					Period:      "年中無休",
					Phone:       "0166-92-1718",
					Price:       "有料。",
					Description: "千代田の丘キャンプ場は、美瑛町のファームズ千代田内に位置し、自然豊かな牧場空間を楽しめます。",
					IconPath:    "/static/img/campsiteflag.jpeg",
				},
			},
			want: struct {
				Period string
				err    error
			}{
				Period: "年中無休",
				err:    nil,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewSpotRepository(db)

			err := repo.Update(tt.in.ctx, tt.in.spot)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.want.err)
			}

			spot, err := repo.GetSpotByID(tt.in.ctx, tt.in.spot.ID.String())
			if err != nil {
				t.Errorf("After Update(), GetSpotByID() error = %v, want no error", err)
			}
			if spot.Period != tt.want.Period {
				t.Errorf("After Update(), GetSpotByID() got name = %v, want name = %v", spot.Period, tt.want.Period)
			}
		})
	}
}

func TestSpotRepo_Delete(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name    string
		setup   func(db *sql.DB)
		in      SpotDeleteArg
		wantErr error
	}{
		{
			name: "Success",
			in: SpotDeleteArg{
				ctx: context.Background(),
				spot: model.Spot{
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
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewSpotRepository(db)

			err := repo.Delete(tt.in.ctx, tt.in.spot.ID.String())

			if err == nil {
				var exists bool
				exists, err = repo.CheckIfSpotExists(tt.in.ctx, tt.in.spot.Lat, tt.in.spot.Lng)
				if exists {
					t.Errorf("Delete() is successful but there is spot in db:  err = %v", err)
				}
			}
		})
	}
}
