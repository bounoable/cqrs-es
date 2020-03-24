package container

import (
	cqrs "github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/aggregate"
	"github.com/bounoable/cqrs-es/command"
)

//go:generate mockgen -source=container.go -destination=../mocks/container/container.go

// Container ...
type Container interface {
	AggregateConfig() aggregate.Config
	SetAggregateConfig(aggregate.Config)
	EventConfig() cqrs.EventConfig
	SetEventConfig(cqrs.EventConfig)
	CommandConfig() command.Config
	SetCommandConfig(command.Config)
	EventBus() cqrs.EventBus
	EventPublisher() cqrs.EventPublisher
	EventSubscriber() cqrs.EventSubscriber
	SetEventBus(cqrs.EventBus)
	EventStore() cqrs.EventStore
	SetEventStore(cqrs.EventStore)
	CommandBus() cqrs.CommandBus
	SetCommandBus(cqrs.CommandBus)
	Snapshots() cqrs.SnapshotRepository
	SetSnapshots(cqrs.SnapshotRepository)
	Aggregates() aggregate.Repository
	SetAggregates(aggregate.Repository)
}

// New ...
func New(aggregates aggregate.Config, events cqrs.EventConfig, commands command.Config) Container {
	return &container{
		aggregateConfig: aggregates,
		eventConfig:     events,
		commandConfig:   commands,
	}
}

type container struct {
	aggregateConfig aggregate.Config
	eventConfig     cqrs.EventConfig
	commandConfig   command.Config
	eventBus        cqrs.EventBus
	eventStore      cqrs.EventStore
	commandBus      cqrs.CommandBus
	snapshots       cqrs.SnapshotRepository
	aggregates      aggregate.Repository
}

func (c *container) AggregateConfig() aggregate.Config {
	return c.aggregateConfig
}

func (c *container) SetAggregateConfig(cfg aggregate.Config) {
	c.aggregateConfig = cfg
}

func (c *container) EventConfig() cqrs.EventConfig {
	return c.eventConfig
}

func (c *container) SetEventConfig(cfg cqrs.EventConfig) {
	c.eventConfig = cfg
}

func (c *container) CommandConfig() command.Config {
	return c.commandConfig
}

func (c *container) SetCommandConfig(cfg command.Config) {
	c.commandConfig = cfg
}

func (c *container) EventBus() cqrs.EventBus {
	return c.eventBus
}

func (c *container) EventPublisher() cqrs.EventPublisher {
	return c.eventBus
}

func (c *container) EventSubscriber() cqrs.EventSubscriber {
	return c.eventBus
}

func (c *container) SetEventBus(cfg cqrs.EventBus) {
	c.eventBus = cfg
}

func (c *container) EventStore() cqrs.EventStore {
	return c.eventStore
}

func (c *container) SetEventStore(store cqrs.EventStore) {
	c.eventStore = store
}

func (c *container) CommandBus() cqrs.CommandBus {
	return c.commandBus
}

func (c *container) SetCommandBus(bus cqrs.CommandBus) {
	c.commandBus = bus
}

func (c *container) Snapshots() cqrs.SnapshotRepository {
	return c.snapshots
}

func (c *container) SetSnapshots(repo cqrs.SnapshotRepository) {
	c.snapshots = repo
}

func (c *container) Aggregates() aggregate.Repository {
	return c.aggregates
}

func (c *container) SetAggregates(repo aggregate.Repository) {
	c.aggregates = repo
}
