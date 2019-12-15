package cqrs_test

import (
	"testing"

	"github.com/bounoable/cqrs"
	mock_cqrs "github.com/bounoable/cqrs/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseAggregate(t *testing.T) {
	typ := cqrs.AggregateType("test")
	id := uuid.New()
	a := cqrs.NewBaseAggregate(typ, id)

	assert.Equal(t, typ, a.AggregateType())
	assert.Equal(t, id, a.AggregateID())
	assert.Equal(t, -1, a.OriginalVersion())
	assert.Equal(t, -1, a.CurrentVersion())
}

func TestTrackChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := cqrs.NewBaseAggregate(cqrs.AggregateType("test"), uuid.New())
	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}
	a.TrackChange(events...)

	assert.Equal(t, events, a.Changes())
	assert.Equal(t, 1, a.CurrentVersion())
}

func TestFlushChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := cqrs.NewBaseAggregate(cqrs.AggregateType("test"), uuid.New())
	events := []cqrs.EventData{
		mock_cqrs.NewMockEventData(ctrl),
		mock_cqrs.NewMockEventData(ctrl),
	}
	a.TrackChange(events...)
	a.FlushChanges()

	assert.Equal(t, 1, a.OriginalVersion())
	assert.Equal(t, 1, a.CurrentVersion())
	assert.Len(t, a.Changes(), 0)
}
