// Code generated by MockGen. DO NOT EDIT.
// Source: commandconfig.go

// Package mock_cqrs is a generated GoMock package.
package mock_cqrs

import (
	cqrs "github.com/bounoable/cqrs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCommandConfig is a mock of CommandConfig interface
type MockCommandConfig struct {
	ctrl     *gomock.Controller
	recorder *MockCommandConfigMockRecorder
}

// MockCommandConfigMockRecorder is the mock recorder for MockCommandConfig
type MockCommandConfigMockRecorder struct {
	mock *MockCommandConfig
}

// NewMockCommandConfig creates a new mock instance
func NewMockCommandConfig(ctrl *gomock.Controller) *MockCommandConfig {
	mock := &MockCommandConfig{ctrl: ctrl}
	mock.recorder = &MockCommandConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandConfig) EXPECT() *MockCommandConfigMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockCommandConfig) Register(arg0 cqrs.CommandType, arg1 cqrs.CommandHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0, arg1)
}

// Register indicates an expected call of Register
func (mr *MockCommandConfigMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockCommandConfig)(nil).Register), arg0, arg1)
}

// Handler mocks base method
func (m *MockCommandConfig) Handler(arg0 cqrs.CommandType) (cqrs.CommandHandler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handler", arg0)
	ret0, _ := ret[0].(cqrs.CommandHandler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handler indicates an expected call of Handler
func (mr *MockCommandConfigMockRecorder) Handler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handler", reflect.TypeOf((*MockCommandConfig)(nil).Handler), arg0)
}

// Handlers mocks base method
func (m *MockCommandConfig) Handlers() map[cqrs.CommandType]cqrs.CommandHandler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handlers")
	ret0, _ := ret[0].(map[cqrs.CommandType]cqrs.CommandHandler)
	return ret0
}

// Handlers indicates an expected call of Handlers
func (mr *MockCommandConfigMockRecorder) Handlers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handlers", reflect.TypeOf((*MockCommandConfig)(nil).Handlers))
}
