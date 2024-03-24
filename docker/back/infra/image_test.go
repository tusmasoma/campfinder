package infra

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type GetSpotImgURLBySpotIDArg struct {
	ctx    context.Context
	spotID string
}

type ImageCreateArg struct {
	ctx context.Context
	img model.Image
}

type ImageDeleteArg struct {
	ctx context.Context
	id  string
}

func TestImageRepo_GetSpotImgURLBySpotID(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    GetSpotImgURLBySpotIDArg
		want  struct {
			imgs []model.Image
			err  error
		}
	}{
		{
			name: "Success",
			in: GetSpotImgURLBySpotIDArg{
				ctx:    context.Background(),
				spotID: "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			},
			want: struct {
				imgs []model.Image
				err  error
			}{
				imgs: []model.Image{
					{
						ID:     uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
						SpotID: uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
						UserID: uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
						URL:    "https://lh3.googleusercontent.com/places/ABCD",
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewImageRepository(db)

			imgs, err := repo.GetSpotImgURLBySpotID(tt.in.ctx, tt.in.spotID)

			ValidateErr(t, err, tt.want.err)

			if d := cmp.Diff(imgs, tt.want.imgs, cmpopts.IgnoreFields(model.Image{}, "Created")); len(d) != 0 {
				t.Errorf("GetSpotImgURLBySpotID() differs: (-got +want)\n%s", d)
			}
		})
	}
}

func TestImageRepo_Create(t *testing.T) {
	patterns := []struct {
		name    string
		setup   func(tx repository.SQLExecutor)
		in      ImageCreateArg
		wantErr error
	}{
		{
			name: "Success",
			in: ImageCreateArg{
				ctx: context.Background(),
				img: model.Image{
					SpotID: uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
					UserID: uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					URL:    "https://lh3.googleusercontent.com/places/XTZ",
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

				repo := NewImageRepository(db)

				err := repo.Create(tt.in.ctx, tt.in.img, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.wantErr)

				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}

func TestImageRepo_Delete(t *testing.T) {
	patterns := []struct {
		name    string
		setup   func(tx repository.SQLExecutor)
		in      ImageDeleteArg
		wantErr error
	}{
		{
			name: "Success",
			in: ImageDeleteArg{
				ctx: context.Background(),
				id:  "31894386-3e60-45a8-bc67-f46b72b42554",
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

				repo := NewImageRepository(db)

				err := repo.Delete(tt.in.ctx, tt.in.id, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.wantErr)

				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}
