package cqrs_test

import (
	"context"
	"testing"

	"github.com/bounoable/cqrs"
	mock_cqrs "github.com/bounoable/cqrs/mocks"
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
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}

	aggregate.EXPECT().AggregateType().Return(aggregateType)
	aggregate.EXPECT().AggregateID().Return(aggregateID)
	aggregate.EXPECT().OriginalVersion().Return(5)
	aggregate.EXPECT().Changes().Return(events)

	eventsArgs := make([]interface{}, len(events))
	for i, e := range events {
		eventsArgs[i] = e
	}

	eventStore.EXPECT().Save(ctx, aggregateType, aggregateID, 5, eventsArgs...).Return(nil)

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
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}

	aggregate.EXPECT().AggregateType().Return(aggregateType)
	aggregate.EXPECT().AggregateID().Return(aggregateID)
	aggregate.EXPECT().OriginalVersion().Return(5)
	aggregate.EXPECT().Changes().Return(events)

	snapshotConfig.EXPECT().IsDue(aggregate).Return(true)
	snapshots.EXPECT().Save(ctx, aggregate).Return(nil)

	eventsArgs := make([]interface{}, len(events))
	for i, e := range events {
		eventsArgs[i] = e
	}

	eventStore.EXPECT().Save(ctx, aggregateType, aggregateID, 5, eventsArgs...).Return(nil)

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

	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}

	eventsArgs := make([]interface{}, len(events))
	for i, e := range events {
		eventsArgs[i] = e
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	aggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().Fetch(ctx, aggregateType, aggregateID, 4, version).Return(events, nil)
	aggregate.EXPECT().ApplyHistory(eventsArgs...).Return(nil)

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

	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}

	eventsArgs := make([]interface{}, len(events))
	for i, e := range events {
		eventsArgs[i] = e
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	snapshots.EXPECT().MaxVersion(ctx, aggregateType, aggregateID, version).Return(snapAggregate, nil)
	snapAggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().Fetch(ctx, aggregateType, aggregateID, 4, version).Return(events, nil)
	snapAggregate.EXPECT().ApplyHistory(eventsArgs...).Return(nil)

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

	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}

	eventsArgs := make([]interface{}, len(events))
	for i, e := range events {
		eventsArgs[i] = e
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	aggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().FetchFrom(ctx, aggregateType, aggregateID, 4).Return(events, nil)
	aggregate.EXPECT().ApplyHistory(eventsArgs...).Return(nil)

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

	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}

	eventsArgs := make([]interface{}, len(events))
	for i, e := range events {
		eventsArgs[i] = e
	}

	aggregateConfig.EXPECT().New(aggregateType, aggregateID).Return(aggregate, nil)
	snapshots.EXPECT().Latest(ctx, aggregateType, aggregateID).Return(snapAggregate, nil)
	snapAggregate.EXPECT().OriginalVersion().Return(3)
	eventStore.EXPECT().FetchFrom(ctx, aggregateType, aggregateID, 4).Return(events, nil)
	snapAggregate.EXPECT().ApplyHistory(eventsArgs...).Return(nil)

	a, err := repo.FetchLatest(ctx, aggregateType, aggregateID)
	assert.Nil(t, err)
	assert.Equal(t, snapAggregate, a)
}
