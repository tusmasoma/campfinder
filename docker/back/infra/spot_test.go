package infra

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
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

			ValidateErr(t, err, tt.want.err)

			if !reflect.DeepEqual(spot, tt.want.spot) {
				t.Errorf("GetSpotByID() \n got = %v,\n want %v", spot, tt.want.spot)
			}
		})
	}
}

func TestSpotRepo_GetSpotByCategory(t *testing.T) {
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

			ValidateErr(t, err, tt.want.err)

			if !reflect.DeepEqual(spots, tt.want.spots) {
				t.Errorf("GetSpotByCategory() \n got = %v,\n want %v", spots, tt.want.spots)
			}
		})
	}
}

func TestSpotRepo_Create(t *testing.T) {
	patterns := []struct {
		name    string
		setup   func(tx repository.SQLExecutor)
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
			txRepo := NewTransactionRepository(db)
			err := txRepo.Transaction(func(tx repository.SQLExecutor) error {
				if tt.setup != nil {
					tt.setup(tx)
				}

				repo := NewSpotRepository(db)

				err := repo.Create(tt.in.ctx, tt.in.spot, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.wantErr)

				if err == nil {
					var exists bool
					exists, err = repo.CheckIfSpotExists(
						tt.in.ctx,
						tt.in.spot.Lat,
						tt.in.spot.Lng,
						repository.QueryOptions{Executor: tx},
					)
					if !exists {
						t.Errorf("Create() is successful but there is not spot in db:  err = %v", err)
					}
				}
				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}

func TestSpotRepo_Update(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(tx repository.SQLExecutor)
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
					ID:          uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
					Category:    "campsite",
					Name:        "旭川市21世紀の森ふれあい広場",
					Address:     "北海道旭川市東旭川町瑞穂4288",
					Lat:         43.7172721,
					Lng:         142.6674615,
					Period:      "年中無休",
					Phone:       "0166-76-2108",
					Price:       "有料。ログハウス大人290円〜750円、高校生以下180〜460円",
					Description: "旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。",
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
			txRepo := NewTransactionRepository(db)
			err := txRepo.Transaction(func(tx repository.SQLExecutor) error {
				if tt.setup != nil {
					tt.setup(tx)
				}

				repo := NewSpotRepository(db)

				err := repo.Update(tt.in.ctx, tt.in.spot, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.want.err)

				spot, err := repo.GetSpotByID(tt.in.ctx, tt.in.spot.ID.String(), repository.QueryOptions{Executor: tx})
				if err != nil {
					t.Errorf("After Update(), GetSpotByID() error = %v, want no error", err)
				}
				if spot.Period != tt.want.Period {
					t.Errorf("After Update(), GetSpotByID() got name = %v, want name = %v", spot.Period, tt.want.Period)
				}
				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}

func TestSpotRepo_Delete(t *testing.T) {
	patterns := []struct {
		name    string
		setup   func(tx repository.SQLExecutor)
		in      SpotDeleteArg
		wantErr error
	}{
		{
			name: "Success",
			in: SpotDeleteArg{
				ctx: context.Background(),
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
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			txRepo := NewTransactionRepository(db)
			err := txRepo.Transaction(func(tx repository.SQLExecutor) error {
				if tt.setup != nil {
					tt.setup(tx)
				}

				repo := NewSpotRepository(db)

				err := repo.Delete(tt.in.ctx, tt.in.spot.ID.String(), repository.QueryOptions{Executor: tx})

				if err == nil {
					var exists bool
					exists, err = repo.CheckIfSpotExists(
						tt.in.ctx,
						tt.in.spot.Lat,
						tt.in.spot.Lng,
						repository.QueryOptions{Executor: tx},
					)
					if exists {
						t.Errorf("Delete() is successful but there is spot in db:  err = %v", err)
					}
				}
				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}
