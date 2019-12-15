// Code generated by MockGen. DO NOT EDIT.
// Source: snapshot.go

// Package mock_cqrs is a generated GoMock package.
package mock_cqrs

import (
	context "context"
	cqrs_es "github.com/bounoable/cqrs-es"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	reflect "reflect"
)

// MockSnapshotRepository is a mock of SnapshotRepository interface
type MockSnapshotRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotRepositoryMockRecorder
}

// MockSnapshotRepositoryMockRecorder is the mock recorder for MockSnapshotRepository
type MockSnapshotRepositoryMockRecorder struct {
	mock *MockSnapshotRepository
}

// NewMockSnapshotRepository creates a new mock instance
func NewMockSnapshotRepository(ctrl *gomock.Controller) *MockSnapshotRepository {
	mock := &MockSnapshotRepository{ctrl: ctrl}
	mock.recorder = &MockSnapshotRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSnapshotRepository) EXPECT() *MockSnapshotRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockSnapshotRepository) Save(ctx context.Context, snap cqrs_es.Aggregate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, snap)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockSnapshotRepositoryMockRecorder) Save(ctx, snap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSnapshotRepository)(nil).Save), ctx, snap)
}

// Find mocks base method
func (m *MockSnapshotRepository) Find(ctx context.Context, typ cqrs_es.AggregateType, id uuid.UUID, version int) (cqrs_es.Aggregate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, typ, id, version)
	ret0, _ := ret[0].(cqrs_es.Aggregate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockSnapshotRepositoryMockRecorder) Find(ctx, typ, id, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockSnapshotRepository)(nil).Find), ctx, typ, id, version)
}

// Latest mocks base method
func (m *MockSnapshotRepository) Latest(ctx context.Context, typ cqrs_es.AggregateType, id uuid.UUID) (cqrs_es.Aggregate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Latest", ctx, typ, id)
	ret0, _ := ret[0].(cqrs_es.Aggregate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Latest indicates an expected call of Latest
func (mr *MockSnapshotRepositoryMockRecorder) Latest(ctx, typ, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Latest", reflect.TypeOf((*MockSnapshotRepository)(nil).Latest), ctx, typ, id)
}

// MaxVersion mocks base method
func (m *MockSnapshotRepository) MaxVersion(ctx context.Context, typ cqrs_es.AggregateType, id uuid.UUID, maxVersion int) (cqrs_es.Aggregate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MaxVersion", ctx, typ, id, maxVersion)
	ret0, _ := ret[0].(cqrs_es.Aggregate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MaxVersion indicates an expected call of MaxVersion
func (mr *MockSnapshotRepositoryMockRecorder) MaxVersion(ctx, typ, id, maxVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaxVersion", reflect.TypeOf((*MockSnapshotRepository)(nil).MaxVersion), ctx, typ, id, maxVersion)
}

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
