// Code generated by MockGen. DO NOT EDIT.
// Source: snapshot.go

// Package mock_aggregate is a generated GoMock package.
package mock_aggregate

import (
	cqrs_es "github.com/bounoable/cqrs-es"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSnapshotConfig is a mock of SnapshotConfig interface
type MockSnapshotConfig struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotConfigMockRecorder
}

// MockSnapshotConfigMockRecorder is the mock recorder for MockSnapshotConfig
type MockSnapshotConfigMockRecorder struct {
	mock *MockSnapshotConfig
}

// NewMockSnapshotConfig creates a new mock instance
func NewMockSnapshotConfig(ctrl *gomock.Controller) *MockSnapshotConfig {
	mock := &MockSnapshotConfig{ctrl: ctrl}
	mock.recorder = &MockSnapshotConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSnapshotConfig) EXPECT() *MockSnapshotConfigMockRecorder {
	return m.recorder
}

// IsDue mocks base method
func (m *MockSnapshotConfig) IsDue(arg0 cqrs_es.Aggregate) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDue", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDue indicates an expected call of IsDue
func (mr *MockSnapshotConfigMockRecorder) IsDue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDue", reflect.TypeOf((*MockSnapshotConfig)(nil).IsDue), arg0)
}
