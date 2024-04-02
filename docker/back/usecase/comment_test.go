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

type CommentCreateArg struct {
	ctx      context.Context
	spotID   uuid.UUID
	starRate float64
	text     string
	user     model.User
}

type CommentUpdateArg struct {
	ctx      context.Context
	id       uuid.UUID
	spotID   uuid.UUID
	userID   uuid.UUID
	starRate float64
	text     string
	user     model.User
}

type CommentDeleteArg struct {
	ctx    context.Context
	id     string
	userID string
	user   model.User
}

func TestCommentUseCase_GetCommentBySpotID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentRepository,
		)
		arg struct {
			ctx    context.Context
			spotID string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentRepository) {
				const layout = "2006-01-02T15:04:05Z"
				created, _ := time.Parse(layout, "0001-01-01T00:00:00Z")
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "SpotID", Value: "fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"}},
				).Return(
					[]model.Comment{
						{
							ID:       uuid.New(),
							SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
							UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
							StarRate: 5.0,
							Text:     "いいスポットでした！!!",
							Created:  created,
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
			repo := mock.NewMockCommentRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			usecase := NewCommentUseCase(repo)

			_, err := usecase.ListComments(tt.arg.ctx, tt.arg.spotID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetCommentBySpotID() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetCommentBySpotID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommentUseCase_CommentCreate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentRepository,
		)
		arg     CommentCreateArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentRepository) {
				comment := model.Comment{
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5.0,
					Text:     "いいスポットでした！!!",
				}
				m.EXPECT().Create(gomock.Any(), comment).Return(nil)
			},
			arg: CommentCreateArg{
				ctx:      context.Background(),
				spotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
				starRate: 5.0,
				text:     "いいスポットでした！!!",
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
			repo := mock.NewMockCommentRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			usecase := NewCommentUseCase(repo)

			err := usecase.CreateComment(tt.arg.ctx, tt.arg.spotID, tt.arg.starRate, tt.arg.text, tt.arg.user)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CommentCreate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CommentCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommentUseCase_CommentUpdate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentRepository,
		)
		arg     CommentUpdateArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentRepository) {
				comment := model.Comment{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5.0,
					Text:     "いいスポットでした！!!",
				}
				m.EXPECT().Update(gomock.Any(), "31894386-3e60-45a8-bc67-f46b72b42554", comment).Return(nil)
			},
			arg: CommentUpdateArg{
				ctx:      context.Background(),
				id:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
				spotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
				userID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
				starRate: 5.0,
				text:     "いいスポットでした！!!",
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
			setup: func(m *mock.MockCommentRepository) {
				comment := model.Comment{
					ID:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
					SpotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
					UserID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
					StarRate: 5.0,
					Text:     "いいスポットでした！!!",
				}
				m.EXPECT().Update(gomock.Any(), "31894386-3e60-45a8-bc67-f46b72b42554", comment).Return(nil)
			},
			arg: CommentUpdateArg{
				ctx:      context.Background(),
				id:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
				spotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
				userID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
				starRate: 5.0,
				text:     "いいスポットでした！!!",
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
			arg: CommentUpdateArg{
				ctx:      context.Background(),
				id:       uuid.MustParse("31894386-3e60-45a8-bc67-f46b72b42554"),
				spotID:   uuid.MustParse("fb816fc7-ddcf-4fa0-9be0-d1fd0b8b5052"),
				userID:   uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"),
				starRate: 5.0,
				text:     "いいスポットでした！!!",
				user: model.User{
					ID:       uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e61234"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				},
			},
			wantErr: fmt.Errorf("don't have permission to update comment"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := mock.NewMockCommentRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			usecase := NewCommentUseCase(repo)

			err := usecase.UpdateComment(
				tt.arg.ctx,
				tt.arg.id,
				tt.arg.spotID,
				tt.arg.userID,
				tt.arg.starRate,
				tt.arg.text,
				tt.arg.user,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CommentUpdate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CommentUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommentUseCase_CommentDelete(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCommentRepository,
		)
		arg     CommentDeleteArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockCommentRepository) {
				m.EXPECT().Delete(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
				).Return(nil)
			},
			arg: CommentDeleteArg{
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
			setup: func(m *mock.MockCommentRepository) {
				m.EXPECT().Delete(
					gomock.Any(),
					"31894386-3e60-45a8-bc67-f46b72b42554",
				).Return(nil)
			},
			arg: CommentDeleteArg{
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
			name: "Fail: Not authorized to delete",
			arg: CommentDeleteArg{
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
			wantErr: fmt.Errorf("don't have permission to delete comment"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := mock.NewMockCommentRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			usecase := NewCommentUseCase(repo)

			err := usecase.DeleteComment(
				tt.arg.ctx,
				tt.arg.id,
				tt.arg.userID,
				tt.arg.user,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CommentDelete() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CommentDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
