package aggregate

import (
	"context"
	"fmt"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/event"
	"github.com/google/uuid"
)

type repository struct {
	eventStore   cqrs.EventStore
	aggregateCfg cqrs.AggregateConfig
	snapshots    cqrs.SnapshotRepository
}

// Repository ...
func Repository(eventStore cqrs.EventStore, aggregateCfg cqrs.AggregateConfig) cqrs.AggregateRepository {
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

	if aggregate.CurrentVersion() != version {
		return nil, UnreachedVersionError{
			Wanted:  version,
			Current: aggregate.CurrentVersion(),
		}
	}

	return aggregate, nil
}

func (r repository) FetchLatest(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error) {
	aggregate, err := r.aggregateCfg.New(typ, id)
	if err != nil {
		return nil, err
	}

	return r.FetchLatestWithBase(ctx, aggregate)
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
	return r.eventStore.RemoveAggregate(ctx, aggregate.AggregateType(), aggregate.AggregateID())
}

func (r repository) RemoveType(ctx context.Context, typ cqrs.AggregateType) error {
	return r.eventStore.RemoveAggregateType(ctx, typ)
}

func (r repository) Query(ctx context.Context, query cqrs.AggregateQuery) (cqrs.AggregateCursor, error) {
	eventQueryOpts := []event.QueryOption{event.QueryVersions(0)}

	if len(query.Types()) > 0 {
		eventQueryOpts = append(eventQueryOpts, event.QueryAggregateType(query.Types()...))
	}

	if len(query.IDs()) > 0 {
		eventQueryOpts = append(eventQueryOpts, event.QueryAggregateID(query.IDs()...))
	}

	eventCur, err := r.eventStore.Query(ctx, event.Query(eventQueryOpts...))
	if err != nil {
		return nil, err
	}

	return newCursor(r, eventCur), nil
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

// UnreachedVersionError ...
type UnreachedVersionError struct {
	Wanted  int
	Current int
}

func (err UnreachedVersionError) Error() string {
	return fmt.Sprintf("want aggregate version %d but current version is %d", err.Wanted, err.Current)
}

type snapshotter struct {
	cqrs.AggregateRepository
	cfg       SnapshotConfig
	snapshots cqrs.SnapshotRepository
}

// Snapshotter ...
func Snapshotter(aggregates cqrs.AggregateRepository, cfg SnapshotConfig, snapshots cqrs.SnapshotRepository) cqrs.AggregateRepository {
	if aggregates == nil {
		panic("nil aggregate repository")
	}

	if cfg == nil {
		panic("nil snapshot config")
	}

	if snapshots == nil {
		panic("nil snapshot repository")
	}

	return snapshotter{
		AggregateRepository: aggregates,
		cfg:                 cfg,
		snapshots:           snapshots,
	}
}

func (s snapshotter) Save(ctx context.Context, aggregate cqrs.Aggregate) error {
	if err := s.AggregateRepository.Save(ctx, aggregate); err != nil {
		return err
	}

	if s.cfg.IsDue(aggregate) {
		return s.snapshots.Save(ctx, aggregate)
	}

	return nil
}

func (s snapshotter) Fetch(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, version int) (cqrs.Aggregate, error) {
	aggregate, err := s.snapshots.MaxVersion(ctx, typ, id, version)
	if err != nil {
		if aggregate, err = s.AggregateRepository.Fetch(ctx, typ, id, version); err != nil {
			return nil, err
		}
	}

	return s.FetchWithBase(ctx, aggregate, version)
}

func (s snapshotter) FetchLatest(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error) {
	aggregate, err := s.snapshots.Latest(ctx, typ, id)
	if err != nil {
		return s.AggregateRepository.FetchLatest(ctx, typ, id)
	}

	return s.FetchLatestWithBase(ctx, aggregate)
}

func (s snapshotter) Remove(ctx context.Context, aggregate cqrs.Aggregate) error {
	if err := s.AggregateRepository.Remove(ctx, aggregate); err != nil {
		return err
	}

	return s.snapshots.RemoveAll(ctx, aggregate.AggregateType(), aggregate.AggregateID())
}
