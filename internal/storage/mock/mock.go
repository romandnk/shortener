// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go
//
// Generated by this command:
//
//	mockgen -source=storage.go -destination=mock/mock.go storage
//
// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	entity "github.com/romandnk/shortener/internal/entity"
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

// CreateURL mocks base method.
func (m *MockURL) CreateURL(ctx context.Context, url entity.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURL", ctx, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateURL indicates an expected call of CreateURL.
func (mr *MockURLMockRecorder) CreateURL(ctx, url any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURL", reflect.TypeOf((*MockURL)(nil).CreateURL), ctx, url)
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