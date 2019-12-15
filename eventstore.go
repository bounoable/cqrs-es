package cqrs

//go:generate mockgen -source=eventstore.go -destination=./mocks/eventstore.go

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// EventStore stores events in a database.
type EventStore interface {
	Save(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, originalVersion int, events ...EventData) error
	Find(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, version int) (EventData, error)
	Fetch(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, from int, to int) ([]EventData, error)
	FetchAll(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID) ([]EventData, error)
	FetchFrom(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, from int) ([]EventData, error)
	FetchTo(ctx context.Context, aggregateType AggregateType, aggregateID uuid.UUID, to int) ([]EventData, error)
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
	LatestVersion   int
	ProvidedVersion int
}

func (err OptimisticConcurrencyError) Error() string {
	return fmt.Sprintf("optimistic concurrency: latest version is %v, provided version is %v", err.LatestVersion, err.ProvidedVersion)
}
