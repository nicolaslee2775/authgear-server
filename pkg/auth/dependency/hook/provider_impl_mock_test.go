// Code generated by MockGen. DO NOT EDIT.
// Source: provider_impl.go

// Package hook is a generated GoMock package.
package hook

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/skygeario/skygear-server/pkg/auth/model"
	reflect "reflect"
)

// MockUserProvider is a mock of UserProvider interface
type MockUserProvider struct {
	ctrl     *gomock.Controller
	recorder *MockUserProviderMockRecorder
}

// MockUserProviderMockRecorder is the mock recorder for MockUserProvider
type MockUserProviderMockRecorder struct {
	mock *MockUserProvider
}

// NewMockUserProvider creates a new mock instance
func NewMockUserProvider(ctrl *gomock.Controller) *MockUserProvider {
	mock := &MockUserProvider{ctrl: ctrl}
	mock.recorder = &MockUserProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserProvider) EXPECT() *MockUserProviderMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockUserProvider) Get(id string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockUserProviderMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserProvider)(nil).Get), id)
}

// UpdateMetadata mocks base method
func (m *MockUserProvider) UpdateMetadata(user *model.User, metadata map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetadata", user, metadata)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetadata indicates an expected call of UpdateMetadata
func (mr *MockUserProviderMockRecorder) UpdateMetadata(user, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetadata", reflect.TypeOf((*MockUserProvider)(nil).UpdateMetadata), user, metadata)
}