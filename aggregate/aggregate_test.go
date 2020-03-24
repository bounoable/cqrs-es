package aggregate_test

import (
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/aggregate"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBase(t *testing.T) {
	typ := cqrs.AggregateType("test")
	id := uuid.New()
	a := aggregate.Base(typ, id)

	assert.Equal(t, typ, a.AggregateType())
	assert.Equal(t, id, a.AggregateID())
	assert.Equal(t, -1, a.OriginalVersion())
	assert.Equal(t, -1, a.CurrentVersion())
}

func TestTrackChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := aggregate.Base(cqrs.AggregateType("test"), uuid.New())
	events := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}
	a.TrackChange(events...)

	assert.Equal(t, events, a.Changes())
	assert.Equal(t, 1, a.CurrentVersion())
}

func TestFlushChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := aggregate.Base(cqrs.AggregateType("test"), uuid.New())
	events := []cqrs.Event{
		mock_cqrs.NewMockEvent(ctrl),
		mock_cqrs.NewMockEvent(ctrl),
	}
	a.TrackChange(events...)
	a.FlushChanges()

	assert.Equal(t, 1, a.OriginalVersion())
	assert.Equal(t, 1, a.CurrentVersion())
	assert.Len(t, a.Changes(), 0)
}
