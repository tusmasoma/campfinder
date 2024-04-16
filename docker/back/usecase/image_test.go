package usecase

import (
	"context"
	"fmt"
	"reflect"
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

func TestImageUseCase_ListImages(t *testing.T) {
	t.Parallel()

	const layout = "2006-01-02T15:04:05Z"
	created, _ := time.Parse(layout, "0001-01-01T00:00:00Z")
	images := model.Images{
		{
			ID:      uuid.New(),
			SpotID:  uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
			UserID:  uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
			URL:     "https://hoge.com/hoge",
			Created: created,
		},
	}
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockImageRepository,
			m1 *mock.MockImagesCacheRepository,
		)
		arg struct {
			ctx    context.Context
			spotID string
		}
		want struct {
			images []model.Image
			err    error
		}
	}{
		{
			name: "success",
			setup: func(m *mock.MockImageRepository, m1 *mock.MockImagesCacheRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "SpotID", Value: "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"}},
				).Return(
					images, nil,
				)
				m1.EXPECT().Set(
					gomock.Any(),
					"fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052",
					images,
				).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				spotID string
			}{
				ctx:    context.Background(),
				spotID: "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052",
			},
			want: struct {
				images []model.Image
				err    error
			}{
				images: images,
				err:    nil,
			},
		},
		{
			name: "success: fail to get images from db, but success to get images from masterdata",
			setup: func(m *mock.MockImageRepository, m1 *mock.MockImagesCacheRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "SpotID", Value: "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"}},
				).Return(
					nil, fmt.Errorf("fail to get images from db"),
				)
				m1.EXPECT().Get(
					gomock.Any(),
					"fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052",
				).Return(&images, nil)
			},
			arg: struct {
				ctx    context.Context
				spotID string
			}{
				ctx:    context.Background(),
				spotID: "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052",
			},
			want: struct {
				images []model.Image
				err    error
			}{
				images: images,
				err:    nil,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			ir := mock.NewMockImageRepository(ctrl)
			ic := mock.NewMockImagesCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ir, ic)
			}

			usecase := NewImageUseCase(ir, ic)

			images, err := usecase.ListImages(tt.arg.ctx, tt.arg.spotID)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetSpotImgURLBySpotID() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetSpotImgURLBySpotID() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(images, tt.want.images) {
				t.Errorf("GetSpot() \n got = %v,\n want %v", images, tt.want)
			}
		})
	}
}

func TestImageUseCase_CreateImage(t *testing.T) {
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
			ir := mock.NewMockImageRepository(ctrl)
			ic := mock.NewMockImagesCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ir)
			}

			usecase := NewImageUseCase(ir, ic)

			err := usecase.CreateImage(tt.arg.ctx, tt.arg.spotID, tt.arg.url, tt.arg.user)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ImageCreate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ImageCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImageUseCase_DeleteImage(t *testing.T) {
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
			ir := mock.NewMockImageRepository(ctrl)
			ic := mock.NewMockImagesCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ir)
			}

			usecase := NewImageUseCase(ir, ic)

			err := usecase.DeleteImage(tt.arg.ctx, tt.arg.id, tt.arg.userID, tt.arg.user)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ImageDelete() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ImageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
