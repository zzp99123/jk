// Code generated by MockGen. DO NOT EDIT.
// Source: webook/internal/repository/cache/user.go

// Package cachemocks is a generated GoMock package.
package cachemocks

import (
	context "context"
	domain "goFoundation/webook/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUsersCacheIF is a mock of UsersCacheIF interface.
type MockUsersCacheIF struct {
	ctrl     *gomock.Controller
	recorder *MockUsersCacheIFMockRecorder
}

// MockUsersCacheIFMockRecorder is the mock recorder for MockUsersCacheIF.
type MockUsersCacheIFMockRecorder struct {
	mock *MockUsersCacheIF
}

// NewMockUsersCacheIF creates a new mock instance.
func NewMockUsersCacheIF(ctrl *gomock.Controller) *MockUsersCacheIF {
	mock := &MockUsersCacheIF{ctrl: ctrl}
	mock.recorder = &MockUsersCacheIFMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersCacheIF) EXPECT() *MockUsersCacheIFMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockUsersCacheIF) Get(ctx context.Context, id int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUsersCacheIFMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUsersCacheIF)(nil).Get), ctx, id)
}

// Id mocks base method.
func (m *MockUsersCacheIF) Id(id int64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id", id)
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id.
func (mr *MockUsersCacheIFMockRecorder) Id(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockUsersCacheIF)(nil).Id), id)
}

// Set mocks base method.
func (m *MockUsersCacheIF) Set(ctx context.Context, u domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockUsersCacheIFMockRecorder) Set(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockUsersCacheIF)(nil).Set), ctx, u)
}
