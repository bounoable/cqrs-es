package cqrs

//go:generate mockgen -source=aggregate.go -destination=./mocks/aggregate.go

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrUnregisteredEventType ...
	ErrUnregisteredEventType = errors.New("unregistered event type")
)

// AggregateType is the type of an aggregate.
type AggregateType string

// Aggregate is the aggregate of an event stream.
type Aggregate interface {
	AggregateID() uuid.UUID
	AggregateType() AggregateType
	OriginalVersion() int
	CurrentVersion() int
	Changes() []Event
	ApplyEvents(...Event) error
	ApplyEvent(Event) error
	ApplyHistory(...Event) error
}

// EventApplier ...
type EventApplier interface {
	Apply(Aggregate, Event) error
}

// EventApplierFunc ...
type EventApplierFunc func(Aggregate, Event) error

// BaseAggregate is the base implementation for an aggregate.
type BaseAggregate struct {
	Aggregate
	ID      uuid.UUID
	Type    AggregateType
	Version int
	changes []Event
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
func (a *BaseAggregate) Changes() []Event {
	return a.changes
}

// TrackChange adds applied events to the aggregate.
func (a *BaseAggregate) TrackChange(events ...Event) {
	a.changes = append(a.changes, events...)
}

// FlushChanges ...
func (a *BaseAggregate) FlushChanges() {
	a.Version = a.CurrentVersion()
	a.changes = nil
}

// ApplyHistory ...
func (a *BaseAggregate) ApplyHistory(events ...Event) error {
	if err := a.ApplyEvents(events...); err != nil {
		return err
	}

	a.FlushChanges()
	return nil
}

// ApplyEvents ...
func (a *BaseAggregate) ApplyEvents(events ...Event) error {
	for _, event := range events {
		if err := a.ApplyEvent(event); err != nil {
			return fmt.Errorf("could not apply event: %w", err)
		}

		a.TrackChange(event)
	}

	return nil
}

// NewEventWithTime ...
func (a *BaseAggregate) NewEventWithTime(typ EventType, data EventData, time time.Time) Event {
	return NewAggregateEventWithTime(typ, data, time, a.AggregateType(), a.AggregateID(), a.CurrentVersion()+1)
}

// NewEvent ...
func (a *BaseAggregate) NewEvent(typ EventType, data EventData) Event {
	return NewAggregateEvent(typ, data, a.AggregateType(), a.AggregateID(), a.CurrentVersion()+1)
}
