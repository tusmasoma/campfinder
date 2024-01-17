package auth

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/db"
	dbmock "github.com/tusmasoma/campfinder/db/mock"
)

func TestHandler_FetchUserFromContext(t *testing.T) {
	t.Parallel()
	passward, _ := db.PasswordEncrypt("password123")
	patterns := []struct {
		name  string
		setup func(
			m *dbmock.MockUserRepository,
		)
		in   func() context.Context
		want struct {
			user db.User
			err  error
		}
	}{
		{
			name: "Success",
			setup: func(m *dbmock.MockUserRepository) {
				m.EXPECT().GetUserByID(
					gomock.Any(),
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
				).Return(
					db.User{
						ID:       uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
					}, nil)
			},
			in: func() context.Context {
				ctx := context.WithValue(
					context.Background(),
					ContextUserIDKey,
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed")
				return ctx
			},
			want: struct {
				user db.User
				err  error
			}{
				user: db.User{
					ID:       uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
				},
				err: nil,
			},
		},
		{
			name: "Fail",
			in: func() context.Context {
				ctx := context.Background()
				return ctx
			},
			want: struct {
				user db.User
				err  error
			}{
				user: db.User{},
				err:  fmt.Errorf("user name not found in request context"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		ctrl := gomock.NewController(t)
		mockUserRepo := dbmock.NewMockUserRepository(ctrl)

		if tt.setup != nil {
			tt.setup(mockUserRepo)
		}

		handler := NewAuthHandler(mockUserRepo)
		user, err := handler.FetchUserFromContext(tt.in())

		if (err != nil) != (tt.want.err != nil) {
			t.Errorf("FetchUserFromContext() error = %v, wantErr %v", err, tt.want.err)
		} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
			t.Errorf("FetchUserFromContext() error = %v, wantErr %v", err, tt.want.err)
		}

		if !reflect.DeepEqual(user, tt.want.user) {
			t.Errorf("FetchUserFromContext() got = %v, want %v", user, tt.want.user)
		}
	}
}
