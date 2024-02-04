package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository/mock"
)

type CreateUserAndGenerateTokenArg struct {
	ctx      context.Context
	email    string
	passward string
}

func TestUserUseCase_CreateUserAndGenerateToken(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockCacheRepository,
		)
		arg     CreateUserAndGenerateTokenArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockCacheRepository) {
				t.Setenv("PRIVATE_KEY_PATH", "../.certificate/private_key.pem")
				m.EXPECT().CheckIfUserExists(
					gomock.Any(),
					"test@gmail.com",
				).Return(false, nil)
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Username already exists",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockCacheRepository) {
				m.EXPECT().CheckIfUserExists(
					gomock.Any(),
					"test@gmail.com",
				).Return(true, nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("user with this email already exists"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr)
			}

			usecase := NewUserUseCase(ur, cr)
			jwt, err := usecase.CreateUserAndGenerateToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateUserAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateUserAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}

func TestUserUseCase_LoginAndGenerateToken(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockCacheRepository,
		)
		arg     CreateUserAndGenerateTokenArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockCacheRepository) {
				t.Setenv("PRIVATE_KEY_PATH", "../.certificate/private_key.pem")
				passward, _ := PasswordEncrypt("password123")
				m.EXPECT().GetUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(
					model.User{
						ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
						IsAdmin:  false,
					}, nil,
				)
				m1.EXPECT().Exists(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return(false)
				m1.EXPECT().Set(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					gomock.Any(),
				).Return(nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: already logged in",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockCacheRepository) {
				passward, _ := PasswordEncrypt("password123")
				m.EXPECT().GetUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(
					model.User{
						ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
						IsAdmin:  false,
					}, nil,
				)
				m1.EXPECT().Exists(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return(true)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("user id in cache"),
		},
		{
			name: "Fail: invalid passward",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockCacheRepository) {
				passward, _ := PasswordEncrypt("password456")
				m.EXPECT().GetUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(
					model.User{
						ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
						IsAdmin:  false,
					}, nil,
				)
				m1.EXPECT().Exists(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return(false)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr)
			}

			usecase := NewUserUseCase(ur, cr)
			jwt, err := usecase.LoginAndGenerateToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("LoginAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("LoginAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}
