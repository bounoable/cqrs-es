// Code generated by MockGen. DO NOT EDIT.
// Source: commandbus.go

// Package mock_cqrs is a generated GoMock package.
package mock_cqrs

import (
	context "context"
	cqrs "github.com/bounoable/cqrs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCommandBus is a mock of CommandBus interface
type MockCommandBus struct {
	ctrl     *gomock.Controller
	recorder *MockCommandBusMockRecorder
}

// MockCommandBusMockRecorder is the mock recorder for MockCommandBus
type MockCommandBusMockRecorder struct {
	mock *MockCommandBus
}

// NewMockCommandBus creates a new mock instance
func NewMockCommandBus(ctrl *gomock.Controller) *MockCommandBus {
	mock := &MockCommandBus{ctrl: ctrl}
	mock.recorder = &MockCommandBusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandBus) EXPECT() *MockCommandBusMockRecorder {
	return m.recorder
}

// Dispatch mocks base method
func (m *MockCommandBus) Dispatch(arg0 context.Context, arg1 cqrs.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dispatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Dispatch indicates an expected call of Dispatch
func (mr *MockCommandBusMockRecorder) Dispatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockCommandBus)(nil).Dispatch), arg0, arg1)
}

// MockCommandHandler is a mock of CommandHandler interface
type MockCommandHandler struct {
	ctrl     *gomock.Controller
	recorder *MockCommandHandlerMockRecorder
}

// MockCommandHandlerMockRecorder is the mock recorder for MockCommandHandler
type MockCommandHandlerMockRecorder struct {
	mock *MockCommandHandler
}

// NewMockCommandHandler creates a new mock instance
func NewMockCommandHandler(ctrl *gomock.Controller) *MockCommandHandler {
	mock := &MockCommandHandler{ctrl: ctrl}
	mock.recorder = &MockCommandHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandHandler) EXPECT() *MockCommandHandlerMockRecorder {
	return m.recorder
}

// HandleCommand mocks base method
func (m *MockCommandHandler) HandleCommand(arg0 context.Context, arg1 cqrs.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleCommand", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleCommand indicates an expected call of HandleCommand
func (mr *MockCommandHandlerMockRecorder) HandleCommand(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleCommand", reflect.TypeOf((*MockCommandHandler)(nil).HandleCommand), arg0, arg1)
}
