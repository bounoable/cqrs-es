package cqrs

//go:generate mockgen -source=aggregate.go -destination=./mocks/aggregate.go

import (
	"github.com/google/uuid"
)

// AggregateType is the type of an aggregate.
type AggregateType string

// Aggregate is the aggregate of an event stream.
type Aggregate interface {
	AggregateID() uuid.UUID
	AggregateType() AggregateType
	OriginalVersion() int
	CurrentVersion() int
	Changes() []EventData
	ApplyEvents(...EventData) error
	ApplyHistory(...EventData) error
}

// BaseAggregate is the base implementation for an aggregate.
type BaseAggregate struct {
	ID      uuid.UUID
	Type    AggregateType
	Version int
	changes []EventData
}

// NewBaseAggregate returns a new BaseAggregate.
func NewBaseAggregate(typ AggregateType, id uuid.UUID) *BaseAggregate {
	return &BaseAggregate{
		ID:      id,
		Type:    typ,
		Version: -1,
	}
}

// AggregateID returns the aggregate ID.
func (a *BaseAggregate) AggregateID() uuid.UUID {
	return a.ID
}

// AggregateType returns the type name of the aggregate.
func (a *BaseAggregate) AggregateType() AggregateType {
	return a.Type
}

// OriginalVersion returns the original version of the aggregate.
func (a *BaseAggregate) OriginalVersion() int {
	return a.Version
}

// CurrentVersion returns the version of the aggregate after applying the changes.
func (a *BaseAggregate) CurrentVersion() int {
	return a.Version + len(a.changes)
}

// Changes returns the applied events.
func (a *BaseAggregate) Changes() []EventData {
	return a.changes
}

// TrackChange adds applied events to the aggregate.
func (a *BaseAggregate) TrackChange(events ...EventData) {
	a.changes = append(a.changes, events...)
}

// FlushChanges ...
func (a *BaseAggregate) FlushChanges() {
	a.Version = a.CurrentVersion()
	a.changes = nil
}
