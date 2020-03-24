package aggregate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/aggregate"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUse__repositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	version := 5
	expectedErr := errors.New("repo error")

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().Fetch(ctx, cqrs.AggregateType("testagg"), id, version).Return(nil, expectedErr)

	agg, err := aggregate.Use(ctx, repo, "testagg", id, version, func(ctx context.Context, agg cqrs.Aggregate) error {
		t.Fatal("function should not have been called")
		return nil
	})

	assert.True(t, errors.Is(err, expectedErr))
	assert.Nil(t, agg)
}

func TestUse__useError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	version := 5
	mockAggregate := mock_cqrs.NewMockAggregate(ctrl)
	expectedErr := errors.New("use error")

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().Fetch(ctx, cqrs.AggregateType("testagg"), id, version).Return(mockAggregate, nil)

	agg, err := aggregate.Use(ctx, repo, "testagg", id, version, func(ctx context.Context, agg cqrs.Aggregate) error {
		return expectedErr
	})

	assert.True(t, errors.Is(err, expectedErr))
	assert.Equal(t, mockAggregate, agg)
}

func TestUse__saveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	version := 5
	saveErr := errors.New("save error")
	mockAggregate := mock_cqrs.NewMockAggregate(ctrl)

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().Fetch(ctx, cqrs.AggregateType("testagg"), id, version).Return(mockAggregate, nil)
	repo.EXPECT().Save(ctx, mockAggregate).Return(saveErr)

	agg, err := aggregate.Use(ctx, repo, "testagg", id, version, func(ctx context.Context, agg cqrs.Aggregate) error {
		return nil
	})

	assert.True(t, errors.Is(err, saveErr))
	assert.Equal(t, mockAggregate, agg)
}

func TestUse__noError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	version := 5
	mockAggregate := mock_cqrs.NewMockAggregate(ctrl)

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().Fetch(ctx, cqrs.AggregateType("testagg"), id, version).Return(mockAggregate, nil)
	repo.EXPECT().Save(ctx, mockAggregate).Return(nil)

	agg, err := aggregate.Use(ctx, repo, "testagg", id, version, func(ctx context.Context, agg cqrs.Aggregate) error {
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, mockAggregate, agg)
}

func TestUseLatest__repositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	expectedErr := errors.New("repo error")

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().FetchLatest(ctx, cqrs.AggregateType("testagg"), id).Return(nil, expectedErr)

	agg, err := aggregate.UseLatest(ctx, repo, "testagg", id, func(ctx context.Context, agg cqrs.Aggregate) error {
		t.Fatal("function should not have been called")
		return nil
	})

	assert.True(t, errors.Is(err, expectedErr))
	assert.Nil(t, agg)
}

func TestUseLatest__useError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	mockAggregate := mock_cqrs.NewMockAggregate(ctrl)
	expectedErr := errors.New("use error")

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().FetchLatest(ctx, cqrs.AggregateType("testagg"), id).Return(mockAggregate, nil)

	agg, err := aggregate.UseLatest(ctx, repo, "testagg", id, func(ctx context.Context, agg cqrs.Aggregate) error {
		return expectedErr
	})

	assert.True(t, errors.Is(err, expectedErr))
	assert.Equal(t, mockAggregate, agg)
}

func TestUseLatest__saveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	mockAggregate := mock_cqrs.NewMockAggregate(ctrl)
	saveErr := errors.New("save error")

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().FetchLatest(ctx, cqrs.AggregateType("testagg"), id).Return(mockAggregate, nil)
	repo.EXPECT().Save(ctx, mockAggregate).Return(saveErr)

	agg, err := aggregate.UseLatest(ctx, repo, "testagg", id, func(ctx context.Context, agg cqrs.Aggregate) error {
		return nil
	})

	assert.True(t, errors.Is(err, saveErr))
	assert.Equal(t, mockAggregate, agg)
}

func TestUseLatest__noError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	id := uuid.New()
	mockAggregate := mock_cqrs.NewMockAggregate(ctrl)

	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
	repo.EXPECT().FetchLatest(ctx, cqrs.AggregateType("testagg"), id).Return(mockAggregate, nil)
	repo.EXPECT().Save(ctx, mockAggregate).Return(nil)

	agg, err := aggregate.UseLatest(ctx, repo, "testagg", id, func(ctx context.Context, agg cqrs.Aggregate) error {
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, mockAggregate, agg)
}
