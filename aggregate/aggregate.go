package aggregate

import (
	"time"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/event"
	"github.com/google/uuid"
)

// BaseAggregate is the base implementation for an aggregate.
type BaseAggregate struct {
	ID      uuid.UUID
	Type    cqrs.AggregateType
	Version int
	changes []cqrs.Event
}

// NewBase returns a new BaseAggregate.
func NewBase(typ cqrs.AggregateType, id uuid.UUID) *BaseAggregate {
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
func (a *BaseAggregate) AggregateType() cqrs.AggregateType {
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
func (a *BaseAggregate) Changes() []cqrs.Event {
	return a.changes
}

// TrackChange adds applied events to the aggregate.
func (a *BaseAggregate) TrackChange(events ...cqrs.Event) {
	a.changes = append(a.changes, events...)
}

// FlushChanges sets the original version to the current version and clears the changes.
func (a *BaseAggregate) FlushChanges() {
	a.Version = a.CurrentVersion()
	a.changes = nil
}

// NewEventWithTime ...
func (a *BaseAggregate) NewEventWithTime(typ cqrs.EventType, data cqrs.EventData, time time.Time) cqrs.Event {
	return event.NewAggregateEventWithTime(typ, data, time, a.AggregateType(), a.AggregateID(), a.CurrentVersion()+1)
}

// NewEvent ...
func (a *BaseAggregate) NewEvent(typ cqrs.EventType, data cqrs.EventData) cqrs.Event {
	return event.NewAggregateEvent(typ, data, a.AggregateType(), a.AggregateID(), a.CurrentVersion()+1)
}

// ApplyHistory ...
func ApplyHistory(aggregate cqrs.Aggregate, events ...cqrs.Event) error {
	if err := ApplyEvents(aggregate, true, events...); err != nil {
		return err
	}

	aggregate.FlushChanges()

	return nil
}

// ApplyEvents ...
func ApplyEvents(aggregate cqrs.Aggregate, track bool, events ...cqrs.Event) error {
	for _, event := range events {
		if err := ApplyEvent(aggregate, track, event); err != nil {
			return err
		}
	}

	return nil
}

// ApplyEvent ...
func ApplyEvent(aggregate cqrs.Aggregate, track bool, event cqrs.Event) error {
	if err := aggregate.ApplyEvent(event); err != nil {
		return err
	}

	if track {
		aggregate.TrackChange(event)
	}

	return nil
}
