// Code generated by MockGen. DO NOT EDIT.
// Source: webook/internal/service/code.go

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCodeServiceIF is a mock of CodeServiceIF interface.
type MockCodeServiceIF struct {
	ctrl     *gomock.Controller
	recorder *MockCodeServiceIFMockRecorder
}

// MockCodeServiceIFMockRecorder is the mock recorder for MockCodeServiceIF.
type MockCodeServiceIFMockRecorder struct {
	mock *MockCodeServiceIF
}

// NewMockCodeServiceIF creates a new mock instance.
func NewMockCodeServiceIF(ctrl *gomock.Controller) *MockCodeServiceIF {
	mock := &MockCodeServiceIF{ctrl: ctrl}
	mock.recorder = &MockCodeServiceIFMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodeServiceIF) EXPECT() *MockCodeServiceIFMockRecorder {
	return m.recorder
}

// Set mocks base method.
func (m *MockCodeServiceIF) Set(ctx context.Context, biz, phone string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, biz, phone)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCodeServiceIFMockRecorder) Set(ctx, biz, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCodeServiceIF)(nil).Set), ctx, biz, phone)
}

// Verify mocks base method.
func (m *MockCodeServiceIF) Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", ctx, biz, phone, expectedCode)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify.
func (mr *MockCodeServiceIFMockRecorder) Verify(ctx, biz, phone, expectedCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockCodeServiceIF)(nil).Verify), ctx, biz, phone, expectedCode)
}
