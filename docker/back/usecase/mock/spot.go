// Code generated by MockGen. DO NOT EDIT.
// Source: spot.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	model "github.com/tusmasoma/campfinder/docker/back/domain/model"
)

// MockSpotUseCase is a mock of SpotUseCase interface.
type MockSpotUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSpotUseCaseMockRecorder
}

// MockSpotUseCaseMockRecorder is the mock recorder for MockSpotUseCase.
type MockSpotUseCaseMockRecorder struct {
	mock *MockSpotUseCase
}

// NewMockSpotUseCase creates a new mock instance.
func NewMockSpotUseCase(ctrl *gomock.Controller) *MockSpotUseCase {
	mock := &MockSpotUseCase{ctrl: ctrl}
	mock.recorder = &MockSpotUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpotUseCase) EXPECT() *MockSpotUseCaseMockRecorder {
	return m.recorder
}

// SpotCreate mocks base method.
func (m *MockSpotUseCase) SpotCreate(ctx context.Context, category, name, address string, lat, lng float64, period, phone, price, description, iconPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpotCreate", ctx, category, name, address, lat, lng, period, phone, price, description, iconPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// SpotCreate indicates an expected call of SpotCreate.
func (mr *MockSpotUseCaseMockRecorder) SpotCreate(ctx, category, name, address, lat, lng, period, phone, price, description, iconPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpotCreate", reflect.TypeOf((*MockSpotUseCase)(nil).SpotCreate), ctx, category, name, address, lat, lng, period, phone, price, description, iconPath)
}

// SpotGet mocks base method.
func (m *MockSpotUseCase) SpotGet(ctx context.Context, categories []string, spotID string) []model.Spot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpotGet", ctx, categories, spotID)
	ret0, _ := ret[0].([]model.Spot)
	return ret0
}

// SpotGet indicates an expected call of SpotGet.
func (mr *MockSpotUseCaseMockRecorder) SpotGet(ctx, categories, spotID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpotGet", reflect.TypeOf((*MockSpotUseCase)(nil).SpotGet), ctx, categories, spotID)
}
