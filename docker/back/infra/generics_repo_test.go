package infra

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

func Test_List_User(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    struct {
			ctx context.Context
			qcs []repository.QueryCondition
		}
		want struct {
			users []model.User
			err   error
		}
	}{
		{
			name: "success",
			in: struct {
				ctx context.Context
				qcs []repository.QueryCondition
			}{
				ctx: context.Background(),
				qcs: []repository.QueryCondition{
					{
						Field: "email",
						Value: "test@gmail.com",
					},
				},
			},
			want: struct {
				users []model.User
				err   error
			}{
				users: []model.User{
					{
						ID:       uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: "password123",
						IsAdmin:  false,
					},
				},
				err: nil,
			},
		},
		{
			name: "Fail",
			in: struct {
				ctx context.Context
				qcs []repository.QueryCondition
			}{
				ctx: context.Background(),
				qcs: []repository.QueryCondition{
					{
						Field: "email",
						Value: "Unregistered@gmail.com",
					},
				},
			},
			want: struct {
				users []model.User
				err   error
			}{
				users: nil,
				err:   sql.ErrNoRows,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			dialect := goqu.Dialect("mysql")
			repo := NewGenericRepository[model.User](db, &dialect, "User")

			users, err := repo.List(tt.in.ctx, tt.in.qcs)

			ValidateErr(t, err, tt.want.err)
			if !reflect.DeepEqual(users, tt.want.users) {
				t.Errorf("GetSpotByID() \n got = %v,\n want %v", users, tt.want.users)
			}
		})
	}
}

func Test_Get_User(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    struct {
			ctx context.Context
			id  string
		}
		want struct {
			user model.User
			err  error
		}
	}{
		{
			name: "Success",
			in: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  "5fe0e237-6b49-11ee-b686-0242c0a87001",
			},
			want: struct {
				user model.User
				err  error
			}{
				user: model.User{
					ID:       uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				},
				err: nil,
			},
		},
		{
			name: "Fail",
			in: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  "5fe0e237-6b49-11ee-b686-0242c0a00000",
			},
			want: struct {
				user model.User
				err  error
			}{
				user: model.User{},
				err:  sql.ErrNoRows,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			dialect := goqu.Dialect("mysql")
			repo := NewGenericRepository[model.User](db, &dialect, "User")

			user, err := repo.Get(tt.in.ctx, tt.in.id)

			ValidateErr(t, err, tt.want.err)
			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetSpotByID() \n got = %v,\n want %v", user, tt.want.user)
			}
		})
	}
}

func Test_Create_User(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    struct {
			ctx   context.Context
			users []model.User
		}
		wantError error
	}{
		{
			name: "Success",
			in: struct {
				ctx   context.Context
				users []model.User
			}{
				ctx: context.Background(),
				users: []model.User{
					{
						Name:     "new_test1",
						Email:    "new_test1@gmail.com",
						Password: "password123",
					},
					{
						Name:     "new_test2",
						Email:    "new_test2@gmail.com",
						Password: "password123",
					},
				},
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			dialect := goqu.Dialect("mysql")
			repo := NewGenericRepository[model.User](db, &dialect, "User")

			err := repo.Create(tt.in.ctx, tt.in.users)

			ValidateErr(t, err, tt.wantError)
		})
	}
}

func Test_Update_User(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    struct {
			ctx  context.Context
			id   string
			user model.User
		}
		wantError error
	}{
		{
			name: "Success",
			in: struct {
				ctx  context.Context
				id   string
				user model.User
			}{
				ctx: context.Background(),
				id:  "5fe0e237-6b49-11ee-b686-0242c0a87001",
				user: model.User{
					ID:       uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					Name:     "updated_user",
					Email:    "test@gmail.com",
					Password: "password123",
				},
			},
			wantError: nil,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			dialect := goqu.Dialect("mysql")
			repo := NewGenericRepository[model.User](db, &dialect, "User")

			err := repo.Update(tt.in.ctx, tt.in.id, tt.in.user)

			ValidateErr(t, err, tt.wantError)

			user, err := repo.Get(tt.in.ctx, tt.in.id)
			if err != nil {
				t.Errorf("After Update(), Get() error = %v, want no error", err)
			}
			if !reflect.DeepEqual(user, tt.in.user) {
				t.Errorf("Get() \n got = %v,\n want %v", user, tt.in.user)
			}
		})
	}
}

func Test_Delete_User(t *testing.T) {
	t.Skip()
	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    struct {
			ctx context.Context
			id  string
		}
		wantError error
	}{
		{
			name: "Success",
			in: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  "5fe0e237-6b49-11ee-b686-0242c0a87001",
			},
			wantError: nil,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			dialect := goqu.Dialect("mysql")
			repo := NewGenericRepository[model.User](db, &dialect, "User")

			err := repo.Delete(tt.in.ctx, tt.in.id)

			ValidateErr(t, err, tt.wantError)
		})
	}
}
