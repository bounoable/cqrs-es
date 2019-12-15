package cqrs

//go:generate mockgen -source=event.go -destination=./mocks/event.go

import (
	"time"

	"github.com/google/uuid"
)

// EventType is the type of an event.
type EventType string

// EventData is the payload of an event.
type EventData interface {
	EventType() EventType
	EventTime() time.Time
}

// BaseEventData ...
type BaseEventData struct {
	Type EventType `bson:"type"`
	Time time.Time `bson:"time"`
}

// EventType ...
func (d BaseEventData) EventType() EventType {
	return d.Type
}

// EventTime ...
func (d BaseEventData) EventTime() time.Time {
	return d.Time
}

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

// NewEvent creates a new event.
func NewEvent(typ EventType, data EventData, time time.Time) Event {
	return &event{
		typ:     typ,
		data:    data,
		time:    time,
		version: -1,
	}
}

// NewAggregateEvent creates a new event.
func NewAggregateEvent(typ EventType, data EventData, time time.Time, aggregateType AggregateType, aggregateID uuid.UUID, version int) Event {
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
