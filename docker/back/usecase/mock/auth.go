// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	model "github.com/tusmasoma/campfinder/docker/back/domain/model"
)

// MockAuthUseCase is a mock of AuthUseCase interface.
type MockAuthUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUseCaseMockRecorder
}

// MockAuthUseCaseMockRecorder is the mock recorder for MockAuthUseCase.
type MockAuthUseCaseMockRecorder struct {
	mock *MockAuthUseCase
}

// NewMockAuthUseCase creates a new mock instance.
func NewMockAuthUseCase(ctrl *gomock.Controller) *MockAuthUseCase {
	mock := &MockAuthUseCase{ctrl: ctrl}
	mock.recorder = &MockAuthUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUseCase) EXPECT() *MockAuthUseCaseMockRecorder {
	return m.recorder
}

// GetUserFromContext mocks base method.
func (m *MockAuthUseCase) GetUserFromContext(ctx context.Context) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFromContext", ctx)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFromContext indicates an expected call of GetUserFromContext.
func (mr *MockAuthUseCaseMockRecorder) GetUserFromContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFromContext", reflect.TypeOf((*MockAuthUseCase)(nil).GetUserFromContext), ctx)
}
