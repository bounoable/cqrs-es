package cqrs_test

import (
	"context"
	"testing"

	"github.com/bounoable/cqrs-es"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAggregateRepository(t *testing.T) {
	assert.Panics(t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cqrs.NewAggregateRepository(
			nil,
			mock_cqrs.NewMockAggregateConfig(ctrl),
			mock_cqrs.NewMockSnapshotConfig(ctrl),
			mock_cqrs.NewMockSnapshotRepository(ctrl),
		)
	})
	assert.Panics(t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cqrs.NewAggregateRepository(
			mock_cqrs.NewMockEventStore(ctrl),
			nil,
			mock_cqrs.NewMockSnapshotConfig(ctrl),
			mock_cqrs.NewMockSnapshotRepository(ctrl),
		)
	})
	assert.NotPanics(t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cqrs.NewAggregateRepository(
			mock_cqrs.NewMockEventStore(ctrl),
			mock_cqrs.NewMockAggregateConfig(ctrl),
			mock_cqrs.NewMockSnapshotConfig(ctrl),
			mock_cqrs.NewMockSnapshotRepository(ctrl),
		)
	})
}

func TestSaveAggregateWithoutSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventStore := mock_cqrs.NewMockEventStore(ctrl)
	aggregateConfig := mock_cqrs.NewMockAggregateConfig(ctrl)
	repo := cqrs.NewAggregateRepository(eventStore, aggregateConfig, nil, nil)

	ctx := context.Background()

	aggregate := mock_cqrs.NewMockAggregate(ctrl)
	changes := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}

	aggregate.EXPECT().OriginalVersion().Return(5)
	aggregate.EXPECT().Changes().Return(changes)

	eventStore.EXPECT().Save(ctx, 5, gomock.Any()).Return(nil)

	err := repo.Save(ctx, aggregate)
	assert.Nil(t, err)
}

func TestSaveAggregateWithSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventStore := mock_cqrs.NewMockEventStore(ctrl)
	aggregateConfig := mock_cqrs.NewMockAggregateConfig(ctrl)
	snapshotConfig := mock_cqrs.NewMockSnapshotConfig(ctrl)
	snapshots := mock_cqrs.NewMockSnapshotRepository(ctrl)
	repo := cqrs.NewAggregateRepository(eventStore, aggregateConfig, snapshotConfig, snapshots)

	ctx := context.Background()

	aggregate := mock_cqrs.NewMockAggregate(ctrl)
	changes := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}

	aggregate.EXPECT().OriginalVersion().Return(5)
	aggregate.EXPECT().Changes().Return(changes)

	snapshotConfig.EXPECT().IsDue(aggregate).Return(true)
	snapshots.EXPECT().Save(ctx, aggregate).Return(nil)

	eventStore.EXPECT().Save(ctx, 5, gomock.Any()).Return(nil)

	err := repo.Save(ctx, aggregate)
	assert.Nil(t, err)
}

func TestFetchAggregateWithoutSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventStore := mock_cqrs.NewMockEventStore(ctrl)
	aggregateConfig := mock_cqrs.NewMockAggregateConfig(ctrl)
	repo := cqrs.NewAggregateRepository(eventStore, aggregateConfig, nil, nil)

	ctx := context.Background()
	aggregate := mock_cqrs.NewMockAggregate(ctrl)
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	version := 10

	changes := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}

	changesArgs := make([]interface{}, len(changes))
	for i, change := range changes {
		changesArgs[i] = change
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	aggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().Fetch(ctx, aggregateType, aggregateID, 4, version).Return(changes, nil)
	aggregate.EXPECT().ApplyHistory(changesArgs...).Return(nil)

	a, err := repo.Fetch(ctx, aggregateType, aggregateID, version)
	assert.Nil(t, err)
	assert.Equal(t, aggregate, a)
}

func TestFetchAggregateWithSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventStore := mock_cqrs.NewMockEventStore(ctrl)
	aggregateConfig := mock_cqrs.NewMockAggregateConfig(ctrl)
	snapshots := mock_cqrs.NewMockSnapshotRepository(ctrl)
	repo := cqrs.NewAggregateRepository(eventStore, aggregateConfig, nil, snapshots)

	ctx := context.Background()
	aggregate := mock_cqrs.NewMockAggregate(ctrl)
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	snapAggregate := mock_cqrs.NewMockAggregate(ctrl)
	version := 10

	changes := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}

	changesArgs := make([]interface{}, len(changes))
	for i, change := range changes {
		changesArgs[i] = change
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	snapshots.EXPECT().MaxVersion(ctx, aggregateType, aggregateID, version).Return(snapAggregate, nil)
	snapAggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().Fetch(ctx, aggregateType, aggregateID, 4, version).Return(changes, nil)
	snapAggregate.EXPECT().ApplyHistory(changesArgs...).Return(nil)

	a, err := repo.Fetch(ctx, aggregateType, aggregateID, version)
	assert.Nil(t, err)
	assert.Equal(t, snapAggregate, a)
}

func TestFetchLatestAggregateWithoutSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventStore := mock_cqrs.NewMockEventStore(ctrl)
	aggregateConfig := mock_cqrs.NewMockAggregateConfig(ctrl)
	repo := cqrs.NewAggregateRepository(eventStore, aggregateConfig, nil, nil)

	ctx := context.Background()
	aggregate := mock_cqrs.NewMockAggregate(ctrl)
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()

	changes := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}

	changesArgs := make([]interface{}, len(changes))
	for i, change := range changes {
		changesArgs[i] = change
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	aggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().FetchFrom(ctx, aggregateType, aggregateID, 4).Return(changes, nil)
	aggregate.EXPECT().ApplyHistory(changesArgs...).Return(nil)

	a, err := repo.FetchLatest(ctx, aggregateType, aggregateID)
	assert.Nil(t, err)
	assert.Equal(t, aggregate, a)
}

func TestFetchLatestAggregateWithSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventStore := mock_cqrs.NewMockEventStore(ctrl)
	aggregateConfig := mock_cqrs.NewMockAggregateConfig(ctrl)
	snapshots := mock_cqrs.NewMockSnapshotRepository(ctrl)
	repo := cqrs.NewAggregateRepository(eventStore, aggregateConfig, nil, snapshots)

	ctx := context.Background()
	aggregate := mock_cqrs.NewMockAggregate(ctrl)
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	snapAggregate := mock_cqrs.NewMockAggregate(ctrl)

	changes := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}

	changesArgs := make([]interface{}, len(changes))
	for i, change := range changes {
		changesArgs[i] = change
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	snapshots.EXPECT().Latest(ctx, aggregateType, aggregateID).Return(snapAggregate, nil)
	snapAggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().FetchFrom(ctx, aggregateType, aggregateID, 4).Return(changes, nil)
	snapAggregate.EXPECT().ApplyHistory(changesArgs...).Return(nil)

	a, err := repo.FetchLatest(ctx, aggregateType, aggregateID)
	assert.Nil(t, err)
	assert.Equal(t, snapAggregate, a)
}
