package cqrs

//go:generate mockgen -source=aggregate.go -destination=./mocks/aggregate.go

import (
	"github.com/google/uuid"
)

// AggregateType is the type of an aggregate.
type AggregateType string

// Aggregate is the aggregate of an event stream.
type Aggregate interface {
	AggregateID() uuid.UUID
	AggregateType() AggregateType
	OriginalVersion() int
	CurrentVersion() int
	Changes() []Event
	TrackChange(...Event)
	FlushChanges()
	ApplyEvent(Event) error
}

func (t AggregateType) String() string {
	return string(t)
}
