package cqrs

//go:generate mockgen -source=container.go -destination=./mocks/container.go

import (
	"log"
)

// Container ...
type Container interface {
	AggregateConfig() AggregateConfig
	EventConfig() EventConfig
	CommandConfig() CommandConfig
	EventBus() EventBus
	EventStore() EventStore
	CommandBus() CommandBus
	Snapshots() SnapshotRepository
	Aggregates() AggregateRepository
}

type container struct {
	logger          *log.Logger
	aggregateConfig AggregateConfig
	eventConfig     EventConfig
	commandConfig   CommandConfig
	eventBus        EventBus
	eventStore      EventStore
	commandBus      CommandBus
	snapshots       SnapshotRepository
	aggregates      AggregateRepository
}

// New ...
func New(
	logger *log.Logger,
	aggregateConfig AggregateConfig,
	eventConfig EventConfig,
	commandConfig CommandConfig,
	eventBus EventBus,
	eventStore EventStore,
	commandBus CommandBus,
	snapshots SnapshotRepository,
	aggregates AggregateRepository,
) Container {
	return &container{
		logger:          logger,
		aggregateConfig: aggregateConfig,
		eventConfig:     eventConfig,
		commandConfig:   commandConfig,
		eventBus:        eventBus,
		eventStore:      eventStore,
		commandBus:      commandBus,
		snapshots:       snapshots,
		aggregates:      aggregates,
	}
}

func (c *container) SetLogger(logger *log.Logger) {
	c.logger = logger
}

func (c *container) AggregateConfig() AggregateConfig {
	return c.aggregateConfig
}

func (c *container) EventConfig() EventConfig {
	return c.eventConfig
}

func (c *container) CommandConfig() CommandConfig {
	return c.commandConfig
}

func (c *container) EventBus() EventBus {
	return c.eventBus
}

func (c *container) EventStore() EventStore {
	return c.eventStore
}

func (c *container) CommandBus() CommandBus {
	return c.commandBus
}

func (c *container) Snapshots() SnapshotRepository {
	return c.snapshots
}

func (c *container) Aggregates() AggregateRepository {
	return c.aggregates
}
