// Code generated by MockGen. DO NOT EDIT.
// Source: spot.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/tusmasoma/campfinder/docker/back/domain/model"
	repository "github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

// MockSpotRepository is a mock of SpotRepository interface.
type MockSpotRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSpotRepositoryMockRecorder
}

// MockSpotRepositoryMockRecorder is the mock recorder for MockSpotRepository.
type MockSpotRepositoryMockRecorder struct {
	mock *MockSpotRepository
}

// NewMockSpotRepository creates a new mock instance.
func NewMockSpotRepository(ctrl *gomock.Controller) *MockSpotRepository {
	mock := &MockSpotRepository{ctrl: ctrl}
	mock.recorder = &MockSpotRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpotRepository) EXPECT() *MockSpotRepositoryMockRecorder {
	return m.recorder
}

// CheckIfSpotExists mocks base method.
func (m *MockSpotRepository) CheckIfSpotExists(ctx context.Context, lat, lng float64, opts ...repository.QueryOptions) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, lat, lng}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckIfSpotExists", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfSpotExists indicates an expected call of CheckIfSpotExists.
func (mr *MockSpotRepositoryMockRecorder) CheckIfSpotExists(ctx, lat, lng interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, lat, lng}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfSpotExists", reflect.TypeOf((*MockSpotRepository)(nil).CheckIfSpotExists), varargs...)
}

// Create mocks base method.
func (m *MockSpotRepository) Create(ctx context.Context, spot model.Spot, opts ...repository.QueryOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, spot}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSpotRepositoryMockRecorder) Create(ctx, spot interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, spot}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSpotRepository)(nil).Create), varargs...)
}

// Delete mocks base method.
func (m *MockSpotRepository) Delete(ctx context.Context, id string, opts ...repository.QueryOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, id}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSpotRepositoryMockRecorder) Delete(ctx, id interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, id}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSpotRepository)(nil).Delete), varargs...)
}

// GetSpotByCategory mocks base method.
func (m *MockSpotRepository) GetSpotByCategory(ctx context.Context, category string, opts ...repository.QueryOptions) ([]model.Spot, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, category}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSpotByCategory", varargs...)
	ret0, _ := ret[0].([]model.Spot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpotByCategory indicates an expected call of GetSpotByCategory.
func (mr *MockSpotRepositoryMockRecorder) GetSpotByCategory(ctx, category interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, category}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpotByCategory", reflect.TypeOf((*MockSpotRepository)(nil).GetSpotByCategory), varargs...)
}

// GetSpotByID mocks base method.
func (m *MockSpotRepository) GetSpotByID(ctx context.Context, id string, opts ...repository.QueryOptions) (model.Spot, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, id}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSpotByID", varargs...)
	ret0, _ := ret[0].(model.Spot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpotByID indicates an expected call of GetSpotByID.
func (mr *MockSpotRepositoryMockRecorder) GetSpotByID(ctx, id interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, id}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpotByID", reflect.TypeOf((*MockSpotRepository)(nil).GetSpotByID), varargs...)
}

// Update mocks base method.
func (m *MockSpotRepository) Update(ctx context.Context, spot model.Spot, opts ...repository.QueryOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, spot}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSpotRepositoryMockRecorder) Update(ctx, spot interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, spot}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSpotRepository)(nil).Update), varargs...)
}

// UpdateOrCreate mocks base method.
func (m *MockSpotRepository) UpdateOrCreate(ctx context.Context, spot model.Spot, opts ...repository.QueryOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, spot}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateOrCreate", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrCreate indicates an expected call of UpdateOrCreate.
func (mr *MockSpotRepositoryMockRecorder) UpdateOrCreate(ctx, spot interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, spot}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrCreate", reflect.TypeOf((*MockSpotRepository)(nil).UpdateOrCreate), varargs...)
}