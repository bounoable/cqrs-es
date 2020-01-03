// Code generated by MockGen. DO NOT EDIT.
// Source: aggregate.go

// Package mock_cqrs is a generated GoMock package.
package mock_cqrs

import (
	cqrs_es "github.com/bounoable/cqrs-es"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	reflect "reflect"
)

// MockAggregate is a mock of Aggregate interface
type MockAggregate struct {
	ctrl     *gomock.Controller
	recorder *MockAggregateMockRecorder
}

// MockAggregateMockRecorder is the mock recorder for MockAggregate
type MockAggregateMockRecorder struct {
	mock *MockAggregate
}

// NewMockAggregate creates a new mock instance
func NewMockAggregate(ctrl *gomock.Controller) *MockAggregate {
	mock := &MockAggregate{ctrl: ctrl}
	mock.recorder = &MockAggregateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAggregate) EXPECT() *MockAggregateMockRecorder {
	return m.recorder
}

// AggregateID mocks base method
func (m *MockAggregate) AggregateID() uuid.UUID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateID")
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

// AggregateID indicates an expected call of AggregateID
func (mr *MockAggregateMockRecorder) AggregateID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateID", reflect.TypeOf((*MockAggregate)(nil).AggregateID))
}

// AggregateType mocks base method
func (m *MockAggregate) AggregateType() cqrs_es.AggregateType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateType")
	ret0, _ := ret[0].(cqrs_es.AggregateType)
	return ret0
}

// AggregateType indicates an expected call of AggregateType
func (mr *MockAggregateMockRecorder) AggregateType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateType", reflect.TypeOf((*MockAggregate)(nil).AggregateType))
}

// OriginalVersion mocks base method
func (m *MockAggregate) OriginalVersion() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OriginalVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// OriginalVersion indicates an expected call of OriginalVersion
func (mr *MockAggregateMockRecorder) OriginalVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OriginalVersion", reflect.TypeOf((*MockAggregate)(nil).OriginalVersion))
}

// CurrentVersion mocks base method
func (m *MockAggregate) CurrentVersion() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// CurrentVersion indicates an expected call of CurrentVersion
func (mr *MockAggregateMockRecorder) CurrentVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentVersion", reflect.TypeOf((*MockAggregate)(nil).CurrentVersion))
}

// Changes mocks base method
func (m *MockAggregate) Changes() []cqrs_es.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Changes")
	ret0, _ := ret[0].([]cqrs_es.Event)
	return ret0
}

// Changes indicates an expected call of Changes
func (mr *MockAggregateMockRecorder) Changes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changes", reflect.TypeOf((*MockAggregate)(nil).Changes))
}

// ApplyEvents mocks base method
func (m *MockAggregate) ApplyEvents(arg0 ...cqrs_es.Event) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ApplyEvents", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyEvents indicates an expected call of ApplyEvents
func (mr *MockAggregateMockRecorder) ApplyEvents(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyEvents", reflect.TypeOf((*MockAggregate)(nil).ApplyEvents), arg0...)
}

// ApplyEvent mocks base method
func (m *MockAggregate) ApplyEvent(arg0 cqrs_es.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyEvent indicates an expected call of ApplyEvent
func (mr *MockAggregateMockRecorder) ApplyEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyEvent", reflect.TypeOf((*MockAggregate)(nil).ApplyEvent), arg0)
}

// ApplyHistory mocks base method
func (m *MockAggregate) ApplyHistory(arg0 ...cqrs_es.Event) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ApplyHistory", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyHistory indicates an expected call of ApplyHistory
func (mr *MockAggregateMockRecorder) ApplyHistory(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyHistory", reflect.TypeOf((*MockAggregate)(nil).ApplyHistory), arg0...)
}
