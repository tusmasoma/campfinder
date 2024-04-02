package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository/mock"
)

type ImageCreateArg struct {
	ctx    context.Context
	spotID uuid.UUID
	url    string
	user   model.User
}

type ImageDeleteArg struct {
	ctx    context.Context
	id     string
	userID string
	user   model.User
}

func TestImageUseCase_GetSpotImgURLBySpotID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockImageRepository,
		)
		arg struct {
			ctx    context.Context
			spotID string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockImageRepository) {
				const layout = "2006-01-02T15:04:05Z"
				created, _ := time.Parse(layout, "0001-01-01T00:00:00Z")
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "SpotID", Value: "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"}},
				).Return(
					[]model.Image{
						{
							ID:      uuid.New(),
							SpotID:  uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
							UserID:  uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
							URL:     "https://hoge.com/hoge",
							Created: created,
						},
					}, nil,
				)
			},
			arg: struct {
				ctx    context.Context
				spotID string
			}{
				ctx:    context.Background(),
				spotID: "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052",
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := mock.NewMockImageRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			usecase := NewImageUseCase(repo)

			_, err := usecase.GetSpotImgURLBySpotID(tt.arg.ctx, tt.arg.spotID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetSpotImgURLBySpotID() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetSpotImgURLBySpotID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImageUseCase_ImageCreate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockImageRepository,
		)
		arg     ImageCreateArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockImageRepository) {
				img := model.Image{
					SpotID: uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID: uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					URL:    "https://hoge.com/hoge",
				}
				m.EXPECT().Create(
					gomock.Any(),
					img,
				).Return(nil)
			},
			arg: ImageCreateArg{
				ctx:    context.Background(),
				spotID: uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
				url:    "https://hoge.com/hoge",
				user: model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := mock.NewMockImageRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			usecase := NewImageUseCase(repo)

			err := usecase.ImageCreate(tt.arg.ctx, tt.arg.spotID, tt.arg.url, tt.arg.user)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ImageCreate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ImageCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImageUseCase_ImageDelete(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockImageRepository,
		)
		arg     ImageDeleteArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockImageRepository) {
				m.EXPECT().Delete(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
				).Return(nil)
			},
			arg: ImageDeleteArg{
				ctx:    context.Background(),
				id:     "31894386-3e60-45a8-bc67-f46b72b42554",
				userID: "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				user: model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				},
			},
			wantErr: nil,
		},
		{
			name: "success: Super User",
			setup: func(m *mock.MockImageRepository) {
				m.EXPECT().Delete(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
				).Return(nil)
			},
			arg: ImageDeleteArg{
				ctx:    context.Background(),
				id:     "31894386-3e60-45a8-bc67-f46b72b42554",
				userID: "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				user: model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "super_user",
					Email:    "super_user@gmail.com",
					Password: "password123",
					IsAdmin:  true,
				},
			},
			wantErr: nil,
		},
		{
			name: "Fail: Not authorized to update",
			arg: ImageDeleteArg{
				ctx:    context.Background(),
				id:     "31894386-3e60-45a8-bc67-f46b72b42554",
				userID: "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				user: model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				},
			},
			wantErr: fmt.Errorf("don't have permission to delete images"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := mock.NewMockImageRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			usecase := NewImageUseCase(repo)

			err := usecase.ImageDelete(tt.arg.ctx, tt.arg.id, tt.arg.userID, tt.arg.user)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ImageDelete() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ImageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
