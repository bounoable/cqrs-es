package cqrs

//go:generate mockgen -source=aggregaterepo.go -destination=./mocks/aggregaterepo.go

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// AggregateRepository ...
type AggregateRepository interface {
	Save(ctx context.Context, aggregate Aggregate) error
	Fetch(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Aggregate, error)
	FetchLatest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error)
}

// AggregateRepositoryOption ...
type AggregateRepositoryOption func(*aggregateRepository)

type aggregateRepository struct {
	eventStore   EventStore
	aggregateCfg AggregateConfig
	snapshots    SnapshotRepository

	mux              sync.RWMutex
	dueForSnapshot   map[Aggregate]struct{}
	snapshotEveryNth int
}

// SnapshotEveryNth ...
func SnapshotEveryNth(n int) AggregateRepositoryOption {
	return func(repo *aggregateRepository) {
		repo.snapshotEveryNth = n
	}
}

// NewAggregateRepository ...
// snapshots is optional.
func NewAggregateRepository(
	eventStore EventStore,
	aggregateCfg AggregateConfig,
	snapshots SnapshotRepository,
	options ...AggregateRepositoryOption,
) AggregateRepository {
	if eventStore == nil {
		panic("event store cannot be nil")
	}

	if aggregateCfg == nil {
		panic("aggregate config cannot be nil")
	}

	repo := &aggregateRepository{
		eventStore:   eventStore,
		aggregateCfg: aggregateCfg,
		snapshots:    snapshots,
	}

	for _, opt := range options {
		opt(repo)
	}

	return repo
}

func (r *aggregateRepository) Save(ctx context.Context, aggregate Aggregate) error {
	changes := aggregate.Changes()

	if err := r.eventStore.Save(ctx, aggregate.OriginalVersion(), changes...); err != nil {
		return err
	}

	if r.snapshots == nil || !r.snapshotDue(aggregate) {
		return nil
	}

	return r.snapshots.Save(ctx, aggregate)
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

	if err := aggregate.ApplyHistory(events...); err != nil {
		return nil, err
	}

	r.checkDue(aggregate, len(events))

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

	if err := aggregate.ApplyHistory(events...); err != nil {
		return nil, err
	}

	r.checkDue(aggregate, len(events))

	return aggregate, nil
}

func (r *aggregateRepository) checkDue(aggregate Aggregate, events int) {
	if r.snapshotEveryNth == 0 || events < r.snapshotEveryNth {
		return
	}

	r.mux.Lock()
	r.dueForSnapshot[aggregate] = struct{}{}
	r.mux.Unlock()
}

func (r *aggregateRepository) snapshotDue(aggregate Aggregate) bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, due := r.dueForSnapshot[aggregate]

	if due {
		delete(r.dueForSnapshot, aggregate)
	}

	return due
}
