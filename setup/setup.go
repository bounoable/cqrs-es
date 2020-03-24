package setup

import (
	"context"
	"log"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/aggregate"
	"github.com/bounoable/cqrs-es/command"
	"github.com/bounoable/cqrs-es/event"
)

var (
	globalAggregateConfig = aggregate.Config()
	globalEventConfig     = event.Config()
)

// Aggregate registers a cqrs.AggregateType globally.
func Aggregate(typ cqrs.AggregateType, factory cqrs.AggregateFactory) {
	globalAggregateConfig.Register(typ, factory)
}

// Event registers a cqrs.EventType globally.
func Event(typ cqrs.EventType, proto cqrs.EventData) {
	globalEventConfig.Register(typ, proto)
}

// Setup is the setup helper.
type Setup struct {
	logger *log.Logger

	aggregateConfig cqrs.AggregateConfig
	eventConfig     cqrs.EventConfig
	commandConfig   cqrs.CommandConfig
	snapshotConfig  aggregate.SnapshotConfig

	eventStoreFactory    EventStoreFactory
	eventBusFactory      EventBusFactory
	commandBusFactory    CommandBusFactory
	snapshotRepoFactory  SnapshotRepositoryFactory
	aggregateRepoFactory AggregateRepositoryFactory
}

// EventStoreFactory ...
type EventStoreFactory func(context.Context, cqrs.Container) (cqrs.EventStore, error)

// EventBusFactory ...
type EventBusFactory func(context.Context, cqrs.Container) (cqrs.EventBus, error)

// CommandBusFactory ...
type CommandBusFactory func(context.Context, cqrs.Container) (cqrs.CommandBus, error)

// SnapshotRepositoryFactory ...
type SnapshotRepositoryFactory func(context.Context, cqrs.Container) (cqrs.SnapshotRepository, error)

// AggregateRepositoryFactory ...
type AggregateRepositoryFactory func(context.Context, cqrs.Container) (cqrs.AggregateRepository, error)

// Option is a setup option.
type Option func(*Setup)

// WithLogger adds a logger to the setup.
func WithLogger(logger *log.Logger) Option {
	return func(s *Setup) {
		s.logger = logger
	}
}

// Logger adds a logger to the setup.
func Logger(logger *log.Logger) Option {
	return WithLogger(logger)
}

// WithAggregateConfig adds a cqrs.AggregateConfig to the setup.
func WithAggregateConfig(cfg cqrs.AggregateConfig) Option {
	return func(s *Setup) {
		for typ, fac := range cfg.Factories() {
			s.aggregateConfig.Register(typ, fac)
		}
	}
}

// WithEventConfig adds a cqrs.EventConfig to the setup.
func WithEventConfig(cfg cqrs.EventConfig) Option {
	return func(s *Setup) {
		for typ, proto := range cfg.Protos() {
			s.eventConfig.Register(typ, proto)
		}
	}
}

// WithEventStoreFactory adds an EventStoreFactory to the setup.
func WithEventStoreFactory(f EventStoreFactory) Option {
	return func(s *Setup) {
		s.eventStoreFactory = f
	}
}

// WithEventBusFactory adds an EventStoreFactory to the setup.
func WithEventBusFactory(f EventBusFactory) Option {
	return func(s *Setup) {
		s.eventBusFactory = f
	}
}

// WithCommandBusFactory adds an EventStoreFactory to the setup.
func WithCommandBusFactory(f CommandBusFactory) Option {
	return func(s *Setup) {
		s.commandBusFactory = f
	}
}

// WithSnapshotRepositoryFactory adds an EventStoreFactory to the setup.
func WithSnapshotRepositoryFactory(f SnapshotRepositoryFactory) Option {
	return func(s *Setup) {
		s.snapshotRepoFactory = f
	}
}

// WithSnapshotConfig adds a snapshot config to the setup.
func WithSnapshotConfig(cfg aggregate.SnapshotConfig) Option {
	return func(s *Setup) {
		s.snapshotConfig = cfg
	}
}

// WithAggregateRepositoryFactory adds an EventStoreFactory to the setup.
func WithAggregateRepositoryFactory(f AggregateRepositoryFactory) Option {
	return func(s *Setup) {
		s.aggregateRepoFactory = f
	}
}

// New initiates a setup.
func New(opts ...Option) *Setup {
	s := Setup{
		aggregateConfig: baseAggregateConfig(),
		eventConfig:     baseEventConfig(),
		commandConfig:   command.Config(),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func baseAggregateConfig() cqrs.AggregateConfig {
	cfg := aggregate.Config()
	for typ, fac := range globalAggregateConfig.Factories() {
		cfg.Register(typ, fac)
	}
	return cfg
}

func baseEventConfig() cqrs.EventConfig {
	cfg := event.Config()
	for typ, proto := range globalEventConfig.Protos() {
		cfg.Register(typ, proto)
	}
	return cfg
}

// Command configures the cqrs.CommandHandler for multiple cqrs.CommandTypes.
func (s *Setup) Command(h cqrs.CommandHandler, types ...cqrs.CommandType) *Setup {
	for _, typ := range types {
		s.commandConfig.Register(typ, h)
	}
	return s
}

// Container initializes the components and returns a cqrs.Container.
func (s *Setup) Container(ctx context.Context) (cqrs.Container, error) {
	var err error
	var eventBus cqrs.EventBus
	var commandBus cqrs.CommandBus
	var eventStore cqrs.EventStore
	var snapshotRepo cqrs.SnapshotRepository
	var aggregateRepo cqrs.AggregateRepository

	container := cqrs.NewContainer(s.aggregateConfig, s.eventConfig, s.commandConfig)

	if s.eventBusFactory != nil {
		if eventBus, err = s.eventBusFactory(ctx, container); err != nil {
			return nil, err
		}
		container.SetEventBus(eventBus)
	}

	if s.commandBusFactory != nil {
		if commandBus, err = s.commandBusFactory(ctx, container); err != nil {
			return nil, err
		}
		container.SetCommandBus(commandBus)
	}

	if s.eventStoreFactory != nil {
		if eventStore, err = s.eventStoreFactory(ctx, container); err != nil {
			return nil, err
		}
		container.SetEventStore(eventStore)
	}

	if s.snapshotRepoFactory != nil {
		if snapshotRepo, err = s.snapshotRepoFactory(ctx, container); err != nil {
			return nil, err
		}
		container.SetSnapshots(snapshotRepo)
	}

	if s.aggregateRepoFactory != nil {
		if aggregateRepo, err = s.aggregateRepoFactory(ctx, container); err != nil {
			return nil, err
		}
	}

	if aggregateRepo == nil && eventStore != nil {
		aggregateRepo = aggregate.Repository(eventStore, s.aggregateConfig)

		if snapshotRepo != nil && s.snapshotConfig != nil {
			aggregateRepo = aggregate.Snapshotter(aggregateRepo, s.snapshotConfig, snapshotRepo)
		}
	}

	if aggregateRepo != nil {
		container.SetAggregates(aggregateRepo)
	}

	if commandBus == nil {
		commandBus = command.BusWithConfig(s.commandConfig)
	}

	container.SetCommandBus(commandBus)

	return container, nil
}
