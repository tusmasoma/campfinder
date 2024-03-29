// Code generated by MockGen. DO NOT EDIT.
// Source: comment.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"

	model "github.com/tusmasoma/campfinder/docker/back/domain/model"
)

// MockCommentUseCase is a mock of CommentUseCase interface.
type MockCommentUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCommentUseCaseMockRecorder
}

// MockCommentUseCaseMockRecorder is the mock recorder for MockCommentUseCase.
type MockCommentUseCaseMockRecorder struct {
	mock *MockCommentUseCase
}

// NewMockCommentUseCase creates a new mock instance.
func NewMockCommentUseCase(ctrl *gomock.Controller) *MockCommentUseCase {
	mock := &MockCommentUseCase{ctrl: ctrl}
	mock.recorder = &MockCommentUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentUseCase) EXPECT() *MockCommentUseCaseMockRecorder {
	return m.recorder
}

// CommentCreate mocks base method.
func (m *MockCommentUseCase) CommentCreate(ctx context.Context, spotID uuid.UUID, starRate float64, text string, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommentCreate", ctx, spotID, starRate, text, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommentCreate indicates an expected call of CommentCreate.
func (mr *MockCommentUseCaseMockRecorder) CommentCreate(ctx, spotID, starRate, text, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommentCreate", reflect.TypeOf((*MockCommentUseCase)(nil).CommentCreate), ctx, spotID, starRate, text, user)
}

// CommentDelete mocks base method.
func (m *MockCommentUseCase) CommentDelete(ctx context.Context, id, userID string, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommentDelete", ctx, id, userID, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommentDelete indicates an expected call of CommentDelete.
func (mr *MockCommentUseCaseMockRecorder) CommentDelete(ctx, id, userID, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommentDelete", reflect.TypeOf((*MockCommentUseCase)(nil).CommentDelete), ctx, id, userID, user)
}

// CommentUpdate mocks base method.
func (m *MockCommentUseCase) CommentUpdate(ctx context.Context, id, spotID, userID uuid.UUID, starRate float64, text string, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommentUpdate", ctx, id, spotID, userID, starRate, text, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommentUpdate indicates an expected call of CommentUpdate.
func (mr *MockCommentUseCaseMockRecorder) CommentUpdate(ctx, id, spotID, userID, starRate, text, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommentUpdate", reflect.TypeOf((*MockCommentUseCase)(nil).CommentUpdate), ctx, id, spotID, userID, starRate, text, user)
}

// GetCommentBySpotID mocks base method.
func (m *MockCommentUseCase) GetCommentBySpotID(ctx context.Context, spotID string) ([]model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentBySpotID", ctx, spotID)
	ret0, _ := ret[0].([]model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentBySpotID indicates an expected call of GetCommentBySpotID.
func (mr *MockCommentUseCaseMockRecorder) GetCommentBySpotID(ctx, spotID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentBySpotID", reflect.TypeOf((*MockCommentUseCase)(nil).GetCommentBySpotID), ctx, spotID)
}
