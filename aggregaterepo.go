package cqrs

//go:generate mockgen -source=aggregaterepo.go -destination=./mocks/aggregaterepo.go

import (
	"context"

	"github.com/google/uuid"
)

// AggregateRepository ...
type AggregateRepository interface {
	Save(ctx context.Context, aggregate Aggregate) error
	Fetch(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Aggregate, error)
	FetchLatest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error)
}

type aggregateRepository struct {
	eventStore     EventStore
	aggregateCfg   AggregateConfig
	snapshotConfig SnapshotConfig
	snapshots      SnapshotRepository
}

// NewAggregateRepository ...
// snapshotConfig & snapshots are optional
func NewAggregateRepository(
	eventStore EventStore,
	aggregateCfg AggregateConfig,
	snapshotConfig SnapshotConfig,
	snapshots SnapshotRepository,
) AggregateRepository {
	if eventStore == nil {
		panic("event store cannot be nil")
	}

	if aggregateCfg == nil {
		panic("aggregate config cannot be nil")
	}

	return &aggregateRepository{
		eventStore:     eventStore,
		aggregateCfg:   aggregateCfg,
		snapshotConfig: snapshotConfig,
		snapshots:      snapshots,
	}
}

func (r *aggregateRepository) Save(ctx context.Context, aggregate Aggregate) error {
	changes := aggregate.Changes()
	events := make([]Event, len(changes))

	for i, change := range changes {
		events[i] = NewAggregateEvent(
			change.EventType(),
			change,
			change.EventTime(),
			aggregate.AggregateType(),
			aggregate.AggregateID(),
			aggregate.OriginalVersion()+i+1,
		)
	}

	if err := r.eventStore.Save(ctx, aggregate.OriginalVersion(), events...); err != nil {
		return err
	}

	if r.snapshots == nil || r.snapshotConfig == nil {
		return nil
	}

	if r.snapshotConfig.IsDue(aggregate) {
		if err := r.snapshots.Save(ctx, aggregate); err != nil {
			return err
		}
	}

	return nil
}

func (r *aggregateRepository) Fetch(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Aggregate, error) {
	aggregate, err := r.aggregateCfg.New(typ, id)
	if err != nil {
		return nil, err
	}

	if version == -1 {
		return aggregate, nil
	}

	if r.snapshots != nil {
		snapshot, err := r.snapshots.MaxVersion(ctx, typ, id, version)
		if err == nil {
			aggregate = snapshot
		}
	}

	events, err := r.eventStore.Fetch(ctx, typ, id, aggregate.OriginalVersion()+1, version)
	if err != nil {
		return nil, err
	}

	eventData := make([]EventData, len(events))
	for i, e := range events {
		eventData[i] = e.Data()
	}

	if err := aggregate.ApplyHistory(eventData...); err != nil {
		return nil, err
	}

	return aggregate, nil
}

func (r *aggregateRepository) FetchLatest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error) {
	aggregate, err := r.aggregateCfg.New(typ, id)
	if err != nil {
		return nil, err
	}

	if r.snapshots != nil {
		snapshot, err := r.snapshots.Latest(ctx, typ, id)
		if err == nil {
			aggregate = snapshot
		}
	}

	events, err := r.eventStore.FetchFrom(ctx, typ, id, aggregate.OriginalVersion()+1)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return aggregate, nil
	}

	eventData := make([]EventData, len(events))
	for i, e := range events {
		eventData[i] = e.Data()
	}

	if err := aggregate.ApplyHistory(eventData...); err != nil {
		return nil, err
	}

	return aggregate, nil
}
