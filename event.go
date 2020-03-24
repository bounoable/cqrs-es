package cqrs

//go:generate mockgen -source=event.go -destination=./mocks/event.go

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// EventType is the type of an event.
type EventType string

// EventData is the payload of an event.
type EventData interface{}

// Event is an event.
type Event interface {
	Type() EventType
	Data() EventData
	Time() time.Time

	AggregateType() AggregateType
	AggregateID() uuid.UUID
	Version() int
}

func (t EventType) String() string {
	return string(t)
}

// EventPublisher publishes events.
type EventPublisher interface {
	Publish(ctx context.Context, events ...Event) error
}

// EventSubscriber subscribes to events.
type EventSubscriber interface {
	Subscribe(ctx context.Context, types ...EventType) (<-chan Event, error)
}

// EventBus is the event bus.
type EventBus interface {
	EventPublisher
	EventSubscriber
}

// EventConfig is the configuration for the events.
type EventConfig interface {
	Register(EventType, EventData)
	// NewData creates an EventData instance for EventType.
	// The returned object has to be a non-pointer struct.
	NewData(EventType) (EventData, error)
	Protos() map[EventType]EventData
}

// EventMatcher ...
type EventMatcher func(Event) bool

// EventStore stores events in a database.
type EventStore interface {
	Save(ctx context.Context, originalVersion int, events ...Event) error
	Find(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, version int) (Event, error)
	Fetch(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, from int, to int) ([]Event, error)
	FetchAll(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID) ([]Event, error)
	FetchFrom(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, from int) ([]Event, error)
	FetchTo(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, to int) ([]Event, error)
	RemoveAll(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID) error
}
