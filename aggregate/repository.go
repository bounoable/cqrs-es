package aggregate

//go:generate mockgen -source=repository.go -destination=../mocks/aggregate/repository.go

import (
	"context"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// Repository ...
type Repository interface {
	Save(ctx context.Context, aggregate cqrs.Aggregate) error
	Fetch(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, version int) (cqrs.Aggregate, error)
	FetchLatest(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error)
	Remove(ctx context.Context, aggregate cqrs.Aggregate) error
}

type repository struct {
	eventStore   cqrs.EventStore
	aggregateCfg cqrs.AggregateConfig
}

// NewRepository ...
func NewRepository(eventStore cqrs.EventStore, aggregateCfg cqrs.AggregateConfig) Repository {
	if eventStore == nil {
		panic("nil event store")
	}

	if aggregateCfg == nil {
		panic("nil aggregate config")
	}

	repo := repository{
		eventStore:   eventStore,
		aggregateCfg: aggregateCfg,
	}

	return repo
}

func (r repository) Save(ctx context.Context, aggregate cqrs.Aggregate) error {
	changes := aggregate.Changes()

	if err := r.eventStore.Save(ctx, aggregate.OriginalVersion(), changes...); err != nil {
		return err
	}

	aggregate.FlushChanges()

	return nil
}

func (r repository) Fetch(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, version int) (cqrs.Aggregate, error) {
	aggregate, err := r.aggregateCfg.New(typ, id)
	if err != nil {
		return nil, err
	}

	if version == -1 {
		return aggregate, nil
	}

	events, err := r.eventStore.Fetch(ctx, typ, id, aggregate.OriginalVersion()+1, version)
	if err != nil {
		return nil, err
	}

	if err := ApplyHistory(aggregate, events...); err != nil {
		return nil, err
	}

	return aggregate, nil
}

func (r repository) FetchLatest(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error) {
	aggregate, err := r.aggregateCfg.New(typ, id)
	if err != nil {
		return nil, err
	}

	events, err := r.eventStore.FetchFrom(ctx, typ, id, aggregate.OriginalVersion()+1)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return aggregate, nil
	}

	if err := ApplyHistory(aggregate, events...); err != nil {
		return nil, err
	}

	return aggregate, nil
}

func (r repository) Remove(ctx context.Context, aggregate cqrs.Aggregate) error {
	return r.eventStore.RemoveAll(ctx, aggregate.AggregateType(), aggregate.AggregateID())
}
