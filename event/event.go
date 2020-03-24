package event

import (
	"fmt"
	"time"

	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

type event struct {
	typ           cqrs.EventType
	data          cqrs.EventData
	time          time.Time
	aggregateType cqrs.AggregateType
	aggregateID   uuid.UUID
	version       int
}

// New creates a new event with time set to time.Now().
func New(typ cqrs.EventType, data cqrs.EventData) cqrs.Event {
	return NewWithTime(typ, data, time.Now())
}

// NewWithTime creates a new event.
func NewWithTime(typ cqrs.EventType, data cqrs.EventData, time time.Time) cqrs.Event {
	return &event{
		typ:     typ,
		data:    data,
		time:    time,
		version: -1,
	}
}

// NewAggregateEvent creates a new aggregate event with time set to time.Now().
func NewAggregateEvent(typ cqrs.EventType, data cqrs.EventData, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, version int) cqrs.Event {
	return NewAggregateEventWithTime(typ, data, time.Now(), aggregateType, aggregateID, version)
}

// NewAggregateEventWithTime creates a new aggregate event.
func NewAggregateEventWithTime(typ cqrs.EventType, data cqrs.EventData, time time.Time, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, version int) cqrs.Event {
	return &event{
		typ:           typ,
		data:          data,
		time:          time,
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		version:       version,
	}
}

func (e *event) Type() cqrs.EventType {
	return e.typ
}

func (e *event) Data() cqrs.EventData {
	return e.data
}

func (e *event) Time() time.Time {
	return e.time
}

func (e *event) AggregateType() cqrs.AggregateType {
	return e.aggregateType
}

func (e *event) AggregateID() uuid.UUID {
	return e.aggregateID
}

func (e *event) Version() int {
	return e.version
}

// NotFoundError ...
type NotFoundError struct {
	AggregateType cqrs.AggregateType
	AggregateID   uuid.UUID
	Version       int
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("event not found for %s:%s@%d", err.AggregateType, err.AggregateID, err.Version)
}
