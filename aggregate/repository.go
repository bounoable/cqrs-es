package aggregate

import (
	"context"
	"fmt"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

type repository struct {
	eventStore   cqrs.EventStore
	aggregateCfg cqrs.AggregateConfig
}

// NewRepository ...
func NewRepository(eventStore cqrs.EventStore, aggregateCfg cqrs.AggregateConfig) cqrs.AggregateRepository {
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

	return r.FetchWithBase(ctx, aggregate, version)
}

func (r repository) FetchLatest(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error) {
	aggregate, err := r.aggregateCfg.New(typ, id)
	if err != nil {
		return nil, err
	}

	return r.FetchLatestWithBase(ctx, aggregate)
}

func (r repository) FetchWithBase(ctx context.Context, aggregate cqrs.Aggregate, version int) (cqrs.Aggregate, error) {
	if version < aggregate.CurrentVersion() {
		return nil, IllegalVersionError{
			Err: fmt.Errorf("version %d is below the base version %d", version, aggregate.CurrentVersion()),
		}
	}

	if version == aggregate.CurrentVersion() {
		return aggregate, nil
	}

	events, err := r.eventStore.Fetch(ctx, aggregate.AggregateType(), aggregate.AggregateID(), aggregate.CurrentVersion()+1, version)
	if err != nil {
		return nil, err
	}

	if err := ApplyHistory(aggregate, events...); err != nil {
		return nil, err
	}

	return aggregate, nil
}

func (r repository) FetchLatestWithBase(ctx context.Context, aggregate cqrs.Aggregate) (cqrs.Aggregate, error) {
	events, err := r.eventStore.FetchFrom(ctx, aggregate.AggregateType(), aggregate.AggregateID(), aggregate.CurrentVersion()+1)
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

// IllegalVersionError ...
type IllegalVersionError struct {
	Err error
}

func (err IllegalVersionError) Error() string {
	return err.Err.Error()
}

func (err IllegalVersionError) Unwrap() error {
	return err.Err
}
