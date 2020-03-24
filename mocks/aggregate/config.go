// Code generated by MockGen. DO NOT EDIT.
// Source: config.go

// Package mock_aggregate is a generated GoMock package.
package mock_aggregate

import (
	cqrs_es "github.com/bounoable/cqrs-es"
	aggregate "github.com/bounoable/cqrs-es/aggregate"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	reflect "reflect"
)

// MockConfig is a mock of Config interface
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockConfig) Register(arg0 cqrs_es.AggregateType, arg1 aggregate.Factory) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0, arg1)
}

// Register indicates an expected call of Register
func (mr *MockConfigMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockConfig)(nil).Register), arg0, arg1)
}

// New mocks base method
func (m *MockConfig) New(arg0 cqrs_es.AggregateType, arg1 uuid.UUID) (cqrs_es.Aggregate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0, arg1)
	ret0, _ := ret[0].(cqrs_es.Aggregate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New
func (mr *MockConfigMockRecorder) New(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockConfig)(nil).New), arg0, arg1)
}

// Factories mocks base method
func (m *MockConfig) Factories() map[cqrs_es.AggregateType]aggregate.Factory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Factories")
	ret0, _ := ret[0].(map[cqrs_es.AggregateType]aggregate.Factory)
	return ret0
}

// Factories indicates an expected call of Factories
func (mr *MockConfigMockRecorder) Factories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Factories", reflect.TypeOf((*MockConfig)(nil).Factories))
}
