package infra

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type GetCommentBySpotIDArg struct {
	ctx    context.Context
	spotID string
}

type GetCommentByIDArg struct {
	ctx context.Context
	id  string
}

type CommentCreateArg struct {
	ctx     context.Context
	comment model.Comment
}

type CommentUpdateArg struct {
	ctx     context.Context
	comment model.Comment
}

type CommentDeleteArg struct {
	ctx context.Context
	id  string
}

func TestCommentRepo_GetCommentBySpotID(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    GetCommentBySpotIDArg
		want  struct {
			comments []model.Comment
			err      error
		}
	}{
		{
			name: "Success",
			in: GetCommentBySpotIDArg{
				ctx:    context.Background(),
				spotID: "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
			},
			want: struct {
				comments []model.Comment
				err      error
			}{
				comments: []model.Comment{
					{
						ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b45524b27"),
						SpotID:   uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
						UserID:   uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
						StarRate: 4.5,
						Text:     "素晴らしい場所でした！",
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

			repo := NewCommentRepository(db)

			comments, err := repo.GetCommentBySpotID(tt.in.ctx, tt.in.spotID)

			ValidateErr(t, err, tt.want.err)

			if d := cmp.Diff(comments, tt.want.comments, cmpopts.IgnoreFields(model.Comment{}, "Created")); len(d) != 0 {
				t.Errorf("GetCommentBySpotID()differs: (-got +want)\n%s", d)
			}
		})
	}
}

func TestCommentRepo_GetCommentByID(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    GetCommentByIDArg
		want  struct {
			comment model.Comment
			err     error
		}
	}{
		{
			name: "Success",
			in: GetCommentByIDArg{
				ctx: context.Background(),
				id:  "31894386-3e60-45a8-bc67-f46b45524b27",
			},
			want: struct {
				comment model.Comment
				err     error
			}{
				comment: model.Comment{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b45524b27"),
					SpotID:   uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
					UserID:   uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					StarRate: 4.5,
					Text:     "素晴らしい場所でした！",
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

			repo := NewCommentRepository(db)

			comment, err := repo.GetCommentByID(tt.in.ctx, tt.in.id)

			ValidateErr(t, err, tt.want.err)

			if d := cmp.Diff(comment, tt.want.comment, cmpopts.IgnoreFields(model.Comment{}, "Created")); len(d) != 0 {
				t.Errorf("GetCommentByID()differs: (-got +want)\n%s", d)
			}
		})
	}
}

func TestCommentRepo_Create(t *testing.T) {
	patterns := []struct {
		name    string
		setup   func(tx repository.SQLExecutor)
		in      CommentCreateArg
		wantErr error
	}{
		{
			name: "Success",
			in: CommentCreateArg{
				ctx: context.Background(),
				comment: model.Comment{
					SpotID:   uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
					UserID:   uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					StarRate: 5,
					Text:     "久しぶりに行ったけど最高だった！",
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

				repo := NewCommentRepository(db)

				err := repo.Create(tt.in.ctx, tt.in.comment, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.wantErr)

				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}

func TestCommentRepo_Update(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(tx repository.SQLExecutor)
		in    CommentUpdateArg
		want  struct {
			Text string
			err  error
		}
	}{
		{
			name: "Success",
			in: CommentUpdateArg{
				ctx: context.Background(),
				comment: model.Comment{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b45524b27"),
					SpotID:   uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
					UserID:   uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					StarRate: 4.5,
					Text:     "素晴らしい場所でした！また行きます！",
					Created:  time.Time{},
				},
			},
			want: struct {
				Text string
				err  error
			}{
				Text: "素晴らしい場所でした！また行きます！",
				err:  nil,
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

				repo := NewCommentRepository(db)

				err := repo.Update(tt.in.ctx, tt.in.comment, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.want.err)

				comment, err := repo.GetCommentByID(tt.in.ctx, tt.in.comment.ID.String(), repository.QueryOptions{Executor: tx})
				if err != nil {
					t.Errorf("After Update(), GetCommentByID() error = %v, want no error", err)
				}
				if comment.Text != tt.want.Text {
					t.Errorf("After Update(), GetCommentByID() got name = %v, want name = %v", comment.Text, tt.want.Text)
				}
				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}

func TestComentRepo_Delete(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(tx repository.SQLExecutor)
		in    CommentDeleteArg
		want  struct {
			comment model.Comment
			err     error
		}
	}{
		{
			name: "Success",
			in: CommentDeleteArg{
				ctx: context.Background(),
				id:  "31894386-3e60-45a8-bc67-f46b45524b27",
			},
			want: struct {
				comment model.Comment
				err     error
			}{
				comment: model.Comment{},
				err:     nil,
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

				repo := NewCommentRepository(db)

				err := repo.Delete(tt.in.ctx, tt.in.id, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.want.err)

				comment, err := repo.GetCommentByID(tt.in.ctx, tt.in.id, repository.QueryOptions{Executor: tx})
				ValidateErr(t, err, sql.ErrNoRows)

				if !reflect.DeepEqual(comment, tt.want.comment) {
					t.Errorf("Delete() \n got = %v,\n want %v", comment, tt.want.comment)
				}

				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}
