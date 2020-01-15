package cqrs

//go:generate mockgen -source=eventstore.go -destination=./mocks/eventstore.go

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

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

// EventStoreError ...
type EventStoreError struct {
	Err       error
	StoreName string
}

func (err EventStoreError) Error() string {
	return fmt.Sprintf("%s event store: %s", err.StoreName, err.Err)
}

// Unwrap ...
func (err EventStoreError) Unwrap() error {
	return err.Err
}

// OptimisticConcurrencyError ...
type OptimisticConcurrencyError struct {
	AggregateType   AggregateType
	AggregateID     uuid.UUID
	LatestVersion   int
	ProvidedVersion int
}

func (err OptimisticConcurrencyError) Error() string {
	return fmt.Sprintf(
		"optimistic concurrency (%s:%s): latest version is %v, provided version is %v",
		err.AggregateType,
		err.AggregateID,
		err.LatestVersion,
		err.ProvidedVersion,
	)
}
