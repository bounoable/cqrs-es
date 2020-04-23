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
	Find(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Event, error)
	Fetch(ctx context.Context, typ AggregateType, id uuid.UUID, from int, to int) ([]Event, error)
	FetchAll(ctx context.Context, typ AggregateType, id uuid.UUID) ([]Event, error)
	FetchFrom(ctx context.Context, typ AggregateType, id uuid.UUID, from int) ([]Event, error)
	FetchTo(ctx context.Context, typ AggregateType, id uuid.UUID, to int) ([]Event, error)
	RemoveAggregate(ctx context.Context, typ AggregateType, id uuid.UUID) error
	RemoveAggregateType(ctx context.Context, typ AggregateType) error
	Query(ctx context.Context, query EventQuery) (EventCursor, error)
}

// EventQuery ...
type EventQuery interface {
	// EventTypes returns the event types to query for.
	// Should query for all event types if none provided.
	EventTypes() []EventType
	// AggregateTypes returns the aggregate types to query for.
	// Should query for all aggregate types if none provided.
	AggregateTypes() []AggregateType
	// AggregateIDs returns the aggregate IDs to query for.
	// Should query for all IDs if none provided.
	AggregateIDs() []uuid.UUID
	// Versions returns the aggregate versions to query for.
	Versions() []int
	// VersionRanges returns version ranges to query for.
	VersionRanges() [][2]int
}

// EventCursor ...
type EventCursor interface {
	// Next moves the cursor to the next event.
	Next(context.Context) bool
	// Event returns the current event.
	Event() Event
	// Err returns the cursor error.
	Err() error
	// Close closes the cursor.
	Close(context.Context) error
}
