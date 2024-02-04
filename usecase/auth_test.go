package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/config"
	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository/mock"
)

func TestAuthUseCase_FetchUserFromContext(t *testing.T) {
	t.Parallel()
	passward, _ := PasswordEncrypt("password123")
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
		)
		in   func() context.Context
		want struct {
			user *model.User
			err  error
		}
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository) {
				m.EXPECT().GetUserByID(
					gomock.Any(),
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
				).Return(
					model.User{
						ID:       uuid.MustParse("5c5323e9-c78f-4dac-94ef-d34ab5ea8fed"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
					}, nil,
				)
			},
			in: func() context.Context {
				ctx := context.WithValue(
					context.Background(),
					config.ContextUserIDKey,
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed")
				return ctx
			},
			want: struct {
				user *model.User
				err  error
			}{
				user: &model.User{
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
				user *model.User
				err  error
			}{
				user: nil,
				err:  fmt.Errorf("user name not found in request context"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		ctrl := gomock.NewController(t)
		mockUserRepo := mock.NewMockUserRepository(ctrl)

		if tt.setup != nil {
			tt.setup(mockUserRepo)
		}

		usecase := NewAuthUseCase(mockUserRepo)
		user, err := usecase.FetchUserFromContext(tt.in())

		if (err != nil) != (tt.want.err != nil) {
			t.Errorf("FetchUserFromContext() error = %v, wantErr %v", err, tt.want.err)
		} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
			t.Errorf("FetchUserFromContext() error = %v, wantErr %v", err, tt.want.err)
		}

		if !reflect.DeepEqual(user, tt.want.user) {
			t.Errorf("FetchUserFromContext() got = %v, want %v", *user, tt.want.user)
		}
	}
}
