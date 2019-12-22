// Code generated by MockGen. DO NOT EDIT.
// Source: commandbus.go

// Package mock_cqrs is a generated GoMock package.
package mock_cqrs

import (
	context "context"
	cqrs_es "github.com/bounoable/cqrs-es"
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
func (m *MockCommandBus) Dispatch(arg0 context.Context, arg1 cqrs_es.Command) error {
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
