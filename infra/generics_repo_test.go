package infra

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository"
)

func Test_List_User(t *testing.T) {
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
