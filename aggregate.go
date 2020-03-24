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
	FetchLatest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error)
	Remove(ctx context.Context, aggregate Aggregate) error
}

// SnapshotRepository ...
type SnapshotRepository interface {
	Save(ctx context.Context, snap Aggregate) error
	Find(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Aggregate, error)
	Latest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error)
	MaxVersion(ctx context.Context, typ AggregateType, id uuid.UUID, maxVersion int) (Aggregate, error)
}
