package cqrs

//go:generate mockgen -source=aggregate.go -destination=./mocks/aggregate.go

import (
	"context"

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
	Changes() []Event
	TrackChange(...Event)
	FlushChanges()
	ApplyEvent(Event) error
}

func (t AggregateType) String() string {
	return string(t)
}

// AggregateFactory ...
type AggregateFactory func(uuid.UUID) Aggregate

// AggregateConfig ...
type AggregateConfig interface {
	Register(AggregateType, AggregateFactory)
	New(AggregateType, uuid.UUID) (Aggregate, error)
	Factories() map[AggregateType]AggregateFactory
}

// AggregateRepository ...
type AggregateRepository interface {
	Save(ctx context.Context, aggregate Aggregate) error
	Fetch(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Aggregate, error)
	FetchWithBase(ctx context.Context, aggregate Aggregate, version int) (Aggregate, error)
	FetchLatest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error)
	FetchLatestWithBase(ctx context.Context, aggregate Aggregate) (Aggregate, error)
	Remove(ctx context.Context, aggregate Aggregate) error
	RemoveType(ctx context.Context, typ AggregateType) error
	Query(ctx context.Context, query AggregateQuery) (AggregateCursor, error)
}

// SnapshotRepository ...
type SnapshotRepository interface {
	Save(ctx context.Context, snap Aggregate) error
	Find(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Aggregate, error)
	Latest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error)
	MaxVersion(ctx context.Context, typ AggregateType, id uuid.UUID, maxVersion int) (Aggregate, error)
	Remove(ctx context.Context, typ AggregateType, id uuid.UUID, version int) error
	RemoveAll(ctx context.Context, typ AggregateType, id uuid.UUID) error
}

// AggregateQuery ...
type AggregateQuery interface {
	// Types returns the aggregate types to query for.
	Types() []AggregateType
	// IDs returns the aggregate IDs to query for.
	IDs() []uuid.UUID
}

// AggregateCursor ...
type AggregateCursor interface {
	// Next moves the cursor to the next aggregate.
	Next(context.Context) bool
	// Aggregate fetches and returns the current aggregate at it's latest version.
	Aggregate(context.Context) (Aggregate, error)
	// Version fetches and returns the current aggregate at the provided version.
	Version(context.Context, int) (Aggregate, error)
	// Err returns the current error.
	Err() error
	// Close closes the cursor.
	Close(context.Context) error
}
