package cqrs

//go:generate mockgen -source=eventbus.go -destination=./mocks/eventbus.go

import (
	"context"
)

// EventPublisher publishes events.
type EventPublisher interface {
	Publish(ctx context.Context, events ...Event) error
}

// EventSubscriber subscribes to events.
type EventSubscriber interface {
	Subscribe(ctx context.Context, typ EventType) (<-chan Event, error)
}

// EventBus is the event bus.
type EventBus interface {
	EventPublisher
	EventSubscriber
}
