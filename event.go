package cqrs

//go:generate mockgen -source=event.go -destination=./mocks/event.go

import (
	"time"

	"github.com/google/uuid"
)

// EventType is the type of an event.
type EventType string

// EventData is the payload of an event.
type EventData interface{}

// Event is an event.
type Event interface {
	Type() EventType
	Data() EventData
	Time() time.Time

	AggregateType() AggregateType
	AggregateID() uuid.UUID
	Version() int
}

type event struct {
	typ           EventType
	data          EventData
	time          time.Time
	aggregateType AggregateType
	aggregateID   uuid.UUID
	version       int
}

// NewEvent creates a new event with time set to time.Now().
func NewEvent(typ EventType, data EventData) Event {
	return NewEventWithTime(typ, data, time.Now())
}

// NewEventWithTime creates a new event.
func NewEventWithTime(typ EventType, data EventData, time time.Time) Event {
	return &event{
		typ:     typ,
		data:    data,
		time:    time,
		version: -1,
	}
}

// NewAggregateEvent creates a new aggregate event with time set to time.Now().
func NewAggregateEvent(typ EventType, data EventData, aggregateType AggregateType, aggregateID uuid.UUID, version int) Event {
	return NewAggregateEventWithTime(typ, data, time.Now(), aggregateType, aggregateID, version)
}

// NewAggregateEventWithTime creates a new aggregate event.
func NewAggregateEventWithTime(typ EventType, data EventData, time time.Time, aggregateType AggregateType, aggregateID uuid.UUID, version int) Event {
	return &event{
		typ:           typ,
		data:          data,
		time:          time,
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		version:       version,
	}
}

func (e *event) Type() EventType {
	return e.typ
}

func (e *event) Data() EventData {
	return e.data
}

func (e *event) Time() time.Time {
	return e.time
}

func (e *event) AggregateType() AggregateType {
	return e.aggregateType
}

func (e *event) AggregateID() uuid.UUID {
	return e.aggregateID
}

func (e *event) Version() int {
	return e.version
}

func (t EventType) String() string {
	return string(t)
}
