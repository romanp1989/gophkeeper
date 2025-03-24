// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/romanp1989/gophkeeper/internal/server/grpc/handlers (interfaces: ISecretService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/romanp1989/gophkeeper/domain"
)

// MockISecretService is a mock of ISecretService interface.
type MockISecretService struct {
	ctrl     *gomock.Controller
	recorder *MockISecretServiceMockRecorder
}

// MockISecretServiceMockRecorder is the mock recorder for MockISecretService.
type MockISecretServiceMockRecorder struct {
	mock *MockISecretService
}

// NewMockISecretService creates a new mock instance.
func NewMockISecretService(ctrl *gomock.Controller) *MockISecretService {
	mock := &MockISecretService{ctrl: ctrl}
	mock.recorder = &MockISecretServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISecretService) EXPECT() *MockISecretServiceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockISecretService) Add(arg0 context.Context, arg1 *domain.Secret) (*domain.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(*domain.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockISecretServiceMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockISecretService)(nil).Add), arg0, arg1)
}

// Delete mocks base method.
func (m *MockISecretService) Delete(arg0 context.Context, arg1 uint64, arg2 domain.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockISecretServiceMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockISecretService)(nil).Delete), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockISecretService) Get(arg0 context.Context, arg1 uint64, arg2 domain.UserID) (*domain.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*domain.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockISecretServiceMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockISecretService)(nil).Get), arg0, arg1, arg2)
}

// GetUserSecrets mocks base method.
func (m *MockISecretService) GetUserSecrets(arg0 context.Context, arg1 domain.UserID) ([]*domain.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSecrets", arg0, arg1)
	ret0, _ := ret[0].([]*domain.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSecrets indicates an expected call of GetUserSecrets.
func (mr *MockISecretServiceMockRecorder) GetUserSecrets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSecrets", reflect.TypeOf((*MockISecretService)(nil).GetUserSecrets), arg0, arg1)
}

// Update mocks base method.
func (m *MockISecretService) Update(arg0 context.Context, arg1 *domain.Secret) (*domain.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*domain.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockISecretServiceMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockISecretService)(nil).Update), arg0, arg1)
}
