package cqrs_test

import (
	"testing"
	"time"

	"github.com/bounoable/cqrs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	mock_cqrs "github.com/bounoable/cqrs/mocks"
	"github.com/golang/mock/gomock"
)

func TestNewEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	typ := cqrs.EventType("test")
	data := mock_cqrs.NewMockEventData(ctrl)
	time := time.Now()
	e := cqrs.NewEvent(typ, data, time)

	assert.Equal(t, typ, e.EventType())
	assert.Equal(t, data, e.EventData())
	assert.Equal(t, time, e.Time())
	assert.Equal(t, cqrs.AggregateType(""), e.AggregateType())
	assert.Equal(t, uuid.Nil, e.AggregateID())
	assert.Equal(t, -1, e.Version())
}

func TestNewAggregateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	typ := cqrs.EventType("test")
	data := mock_cqrs.NewMockEventData(ctrl)
	time := time.Now()
	aggregateType := cqrs.AggregateType("test")
	aggregateID := uuid.New()
	version := 5
	e := cqrs.NewAggregateEvent(typ, data, time, aggregateType, aggregateID, 5)

	assert.Equal(t, typ, e.EventType())
	assert.Equal(t, data, e.EventData())
	assert.Equal(t, time, e.Time())
	assert.Equal(t, aggregateType, e.AggregateType())
	assert.Equal(t, aggregateID, e.AggregateID())
	assert.Equal(t, version, e.Version())
}
