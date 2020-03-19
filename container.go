package cqrs

//go:generate mockgen -source=container.go -destination=./mocks/container.go

// Container ...
type Container interface {
	AggregateConfig() AggregateConfig
	SetAggregateConfig(AggregateConfig)
	EventConfig() EventConfig
	SetEventConfig(EventConfig)
	CommandConfig() CommandConfig
	SetCommandConfig(CommandConfig)
	EventBus() EventBus
	EventPublisher() EventPublisher
	EventSubscriber() EventSubscriber
	SetEventBus(EventBus)
	EventStore() EventStore
	SetEventStore(EventStore)
	CommandBus() CommandBus
	SetCommandBus(CommandBus)
	Snapshots() SnapshotRepository
	SetSnapshots(SnapshotRepository)
	Aggregates() AggregateRepository
	SetAggregates(AggregateRepository)
}

// NewContainer ...
func NewContainer(aggregates AggregateConfig, events EventConfig, commands CommandConfig) Container {
	return &container{
		aggregateConfig: aggregates,
		eventConfig:     events,
		commandConfig:   commands,
	}
}

type container struct {
	aggregateConfig AggregateConfig
	eventConfig     EventConfig
	commandConfig   CommandConfig
	eventBus        EventBus
	eventStore      EventStore
	commandBus      CommandBus
	snapshots       SnapshotRepository
	aggregates      AggregateRepository
}

func (c *container) AggregateConfig() AggregateConfig {
	return c.aggregateConfig
}

func (c *container) SetAggregateConfig(cfg AggregateConfig) {
	c.aggregateConfig = cfg
}

func (c *container) EventConfig() EventConfig {
	return c.eventConfig
}

func (c *container) SetEventConfig(cfg EventConfig) {
	c.eventConfig = cfg
}

func (c *container) CommandConfig() CommandConfig {
	return c.commandConfig
}

func (c *container) SetCommandConfig(cfg CommandConfig) {
	c.commandConfig = cfg
}

func (c *container) EventBus() EventBus {
	return c.eventBus
}

func (c *container) EventPublisher() EventPublisher {
	return c.eventBus
}

func (c *container) EventSubscriber() EventSubscriber {
	return c.eventBus
}

func (c *container) SetEventBus(cfg EventBus) {
	c.eventBus = cfg
}

func (c *container) EventStore() EventStore {
	return c.eventStore
}

func (c *container) SetEventStore(store EventStore) {
	c.eventStore = store
}

func (c *container) CommandBus() CommandBus {
	return c.commandBus
}

func (c *container) SetCommandBus(bus CommandBus) {
	c.commandBus = bus
}

func (c *container) Snapshots() SnapshotRepository {
	return c.snapshots
}

func (c *container) SetSnapshots(repo SnapshotRepository) {
	c.snapshots = repo
}

func (c *container) Aggregates() AggregateRepository {
	return c.aggregates
}

func (c *container) SetAggregates(repo AggregateRepository) {
	c.aggregates = repo
}
