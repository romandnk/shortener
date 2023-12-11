// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mock/mock.go service
//
// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockURL is a mock of URL interface.
type MockURL struct {
	ctrl     *gomock.Controller
	recorder *MockURLMockRecorder
}

// MockURLMockRecorder is the mock recorder for MockURL.
type MockURLMockRecorder struct {
	mock *MockURL
}

// NewMockURL creates a new mock instance.
func NewMockURL(ctrl *gomock.Controller) *MockURL {
	mock := &MockURL{ctrl: ctrl}
	mock.recorder = &MockURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURL) EXPECT() *MockURLMockRecorder {
	return m.recorder
}

// CreateURLAlias mocks base method.
func (m *MockURL) CreateURLAlias(ctx context.Context, original string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURLAlias", ctx, original)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateURLAlias indicates an expected call of CreateURLAlias.
func (mr *MockURLMockRecorder) CreateURLAlias(ctx, original any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURLAlias", reflect.TypeOf((*MockURL)(nil).CreateURLAlias), ctx, original)
}

// GetOriginalByAlias mocks base method.
func (m *MockURL) GetOriginalByAlias(ctx context.Context, alias string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalByAlias", ctx, alias)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalByAlias indicates an expected call of GetOriginalByAlias.
func (mr *MockURLMockRecorder) GetOriginalByAlias(ctx, alias any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalByAlias", reflect.TypeOf((*MockURL)(nil).GetOriginalByAlias), ctx, alias)
}
