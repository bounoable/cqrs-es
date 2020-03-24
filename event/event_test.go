package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/event"
	mock_cqrs "github.com/bounoable/cqrs-es/mocks"
	"github.com/golang/mock/gomock"
)

func TestNewWithTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	typ := cqrs.EventType("test")
	data := mock_cqrs.NewMockEventData(ctrl)
	time := time.Now()
	e := event.NewWithTime(typ, data, time)

	assert.Equal(t, typ, e.Type())
	assert.Equal(t, data, e.Data())
	assert.Equal(t, time, e.Time())
	assert.Equal(t, cqrs.AggregateType(""), e.AggregateType())
	assert.Equal(t, uuid.Nil, e.AggregateID())
	assert.Equal(t, -1, e.Version())
}

func TestNewAggregateWithTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	typ := cqrs.EventType("test")
	data := mock_cqrs.NewMockEventData(ctrl)
	time := time.Now()
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	version := 5
	e := event.NewAggregateEventWithTime(typ, data, time, aggregateType, aggregateID, 5)

	assert.Equal(t, typ, e.Type())
	assert.Equal(t, data, e.Data())
	assert.Equal(t, time, e.Time())
	assert.Equal(t, aggregateType, e.AggregateType())
	assert.Equal(t, aggregateID, e.AggregateID())
	assert.Equal(t, version, e.Version())
}
