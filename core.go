package cqrs

//go:generate mockgen -source=core.go -destination=./mocks/core.go

import (
	"context"
	"fmt"
	"log"
)

// EventStoreFactory ...
type EventStoreFactory func(context.Context, Core) (EventStore, error)

// EventBusFactory ...
type EventBusFactory func(context.Context, Core) (EventBus, error)

// CommandBusFactory ...
type CommandBusFactory func(context.Context, Core) (CommandBus, error)

// CommandHandlerFactory ...
type CommandHandlerFactory func(context.Context, Core) (CommandHandler, []CommandType, error)

// SnapshotRepositoryFactory ...
type SnapshotRepositoryFactory func(context.Context, Core) (SnapshotRepository, error)

// AggregateRepositoryFactory ...
type AggregateRepositoryFactory func(context.Context, Core) (AggregateRepository, error)

// Core ...
type Core interface {
	AggregateConfig() AggregateConfig
	EventConfig() EventConfig
	CommandConfig() CommandConfig
	EventBus() EventBus
	EventStore() EventStore
	CommandBus() CommandBus
	Snapshots() SnapshotRepository
	Aggregates() AggregateRepository
}

// Option ...
type Option func(Setup)

// Setup ...
type Setup interface {
	SetLogger(*log.Logger)
	SetAggregateConfig(AggregateConfig)
	SetEventConfig(EventConfig)
	SetCommandConfig(CommandConfig)

	SetEventStoreFactory(EventStoreFactory)
	SetEventBusFactory(EventBusFactory)
	SetCommandBusFactory(CommandBusFactory)
	AddCommandHandlerFactory(...CommandHandlerFactory)
	SetSnapshotRepositoryFactory(SnapshotRepositoryFactory)
	SetAggregateRepositoryFactory(AggregateRepositoryFactory)

	RegisterAggregate(AggregateType, AggregateFactory)
	RegisterEvent(EventType, EventDataFactory)
	RegisterCommand(CommandType, CommandHandler)
}

// WithLogger ...
func WithLogger(logger *log.Logger) Option {
	return func(c Setup) {
		c.SetLogger(logger)
	}
}

// WithAggregateConfig ...
func WithAggregateConfig(cfg AggregateConfig) Option {
	return func(c Setup) {
		c.SetAggregateConfig(cfg)
	}
}

// WithEventConfig ...
func WithEventConfig(cfg EventConfig) Option {
	return func(c Setup) {
		c.SetEventConfig(cfg)
	}
}

// WithCommandConfig ...
func WithCommandConfig(cfg CommandConfig) Option {
	return func(c Setup) {
		c.SetCommandConfig(cfg)
	}
}

// WithAggregate ...
func WithAggregate(typ AggregateType, factory AggregateFactory) Option {
	return func(c Setup) {
		c.RegisterAggregate(typ, factory)
	}
}

// WithEvent ...
func WithEvent(typ EventType, factory EventDataFactory) Option {
	return func(c Setup) {
		c.RegisterEvent(typ, factory)
	}
}

// WithCommand ...
func WithCommand(typ CommandType, handler CommandHandler) Option {
	return func(c Setup) {
		c.RegisterCommand(typ, handler)
	}
}

// WithEventStoreFactory ...
func WithEventStoreFactory(f EventStoreFactory) Option {
	return func(setup Setup) {
		setup.SetEventStoreFactory(f)
	}
}

// WithEventBusFactory ...
func WithEventBusFactory(f EventBusFactory) Option {
	return func(setup Setup) {
		setup.SetEventBusFactory(f)
	}
}

// WithCommandBusFactory ...
func WithCommandBusFactory(f CommandBusFactory) Option {
	return func(setup Setup) {
		setup.SetCommandBusFactory(f)
	}
}

// WithCommandHandlerFactory ...
func WithCommandHandlerFactory(factories ...CommandHandlerFactory) Option {
	return func(setup Setup) {
		setup.AddCommandHandlerFactory(factories...)
	}
}

// WithSnapshotRepositoryFactory ...
func WithSnapshotRepositoryFactory(f SnapshotRepositoryFactory) Option {
	return func(setup Setup) {
		setup.SetSnapshotRepositoryFactory(f)
	}
}

// WithAggregateRepositoryFactory ...
func WithAggregateRepositoryFactory(f AggregateRepositoryFactory) Option {
	return func(setup Setup) {
		setup.SetAggregateRepositoryFactory(f)
	}
}

type core struct {
	logger          *log.Logger
	aggregateConfig AggregateConfig
	eventConfig     EventConfig
	commandConfig   CommandConfig
	snapshotConfig  SnapshotConfig
	eventBus        EventBus
	eventStore      EventStore
	commandBus      CommandBus
	snapshots       SnapshotRepository
	aggregates      AggregateRepository
}

// New ...
func New(ctx context.Context, options ...Option) (Core, error) {
	return NewWithConfigs(ctx, NewAggregateConfig(), NewEventConfig(), NewCommandConfig(), options...)
}

// NewWithConfigs ...
func NewWithConfigs(
	ctx context.Context,
	aggregateConfig AggregateConfig,
	eventConfig EventConfig,
	commandConfig CommandConfig,
	options ...Option,
) (Core, error) {
	c := &core{
		aggregateConfig: aggregateConfig,
		eventConfig:     eventConfig,
		commandConfig:   commandConfig,
	}

	setup := newSetup(c)
	for _, opt := range options {
		opt(setup)
	}

	for typ, factory := range setup.aggregateConfig.Factories() {
		c.aggregateConfig.Register(typ, factory)
	}

	for typ, factory := range setup.eventConfig.Factories() {
		c.eventConfig.Register(typ, factory)
	}

	for typ, handler := range setup.commandConfig.Handlers() {
		c.commandConfig.Register(typ, handler)
	}

	if setup.eventStoreFactory != nil {
		store, err := setup.eventStoreFactory(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("cqrs setup: %w", err)
		}
		c.eventStore = store
	} else if c.logger != nil {
		c.logger.Println("cqrs setup: no event store")
	}

	if setup.eventBusFactory != nil {
		bus, err := setup.eventBusFactory(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("cqrs setup: %w", err)
		}
		c.eventBus = bus
	} else if c.logger != nil {
		c.logger.Println("cqrs setup: no event bus")
	}

	if setup.commandBusFactory != nil {
		bus, err := setup.commandBusFactory(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("cqrs setup: %w", err)
		}
		c.commandBus = bus
	} else if c.logger != nil {
		c.logger.Println("cqrs setup: no command bus")
	}

	if setup.snapshotRepoFactory != nil {
		repo, err := setup.snapshotRepoFactory(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("cqrs setup: %w", err)
		}
		c.snapshots = repo
	} else if c.logger != nil {
		c.logger.Println("cqrs setup: no snapshot repository")
	}

	if setup.aggregateRepoFactory != nil {
		repo, err := setup.aggregateRepoFactory(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("cqrs setup: %w", err)
		}
		c.aggregates = repo
	} else if c.eventStore != nil && c.aggregateConfig != nil {
		c.aggregates = NewAggregateRepository(c.eventStore, c.aggregateConfig, c.snapshotConfig, c.snapshots)
	} else if c.logger != nil {
		c.logger.Println("cqrs setup: no aggregate repository")
	}

	for _, f := range setup.commandHandlerFactories {
		handler, types, err := f(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("cqrs setup: %w", err)
		}

		for _, typ := range types {
			c.commandConfig.Register(typ, handler)
		}
	}

	return c, nil
}

func (c *core) SetLogger(logger *log.Logger) {
	c.logger = logger
}

func (c *core) AggregateConfig() AggregateConfig {
	return c.aggregateConfig
}

func (c *core) EventConfig() EventConfig {
	return c.eventConfig
}

func (c *core) CommandConfig() CommandConfig {
	return c.commandConfig
}

func (c *core) EventBus() EventBus {
	return c.eventBus
}

func (c *core) EventStore() EventStore {
	return c.eventStore
}

func (c *core) CommandBus() CommandBus {
	return c.commandBus
}

func (c *core) Snapshots() SnapshotRepository {
	return c.snapshots
}

func (c *core) Aggregates() AggregateRepository {
	return c.aggregates
}

type setup struct {
	*core
	aggregateConfig AggregateConfig
	eventConfig     EventConfig
	commandConfig   CommandConfig

	eventStoreFactory       EventStoreFactory
	eventBusFactory         EventBusFactory
	commandBusFactory       CommandBusFactory
	commandHandlerFactories []CommandHandlerFactory
	snapshotRepoFactory     SnapshotRepositoryFactory
	aggregateRepoFactory    AggregateRepositoryFactory
}

func newSetup(c *core) *setup {
	return &setup{
		core:            c,
		aggregateConfig: NewAggregateConfig(),
		eventConfig:     NewEventConfig(),
		commandConfig:   NewCommandConfig(),
	}
}

func (s *setup) SetAggregateConfig(cfg AggregateConfig) {
	s.aggregateConfig = cfg
}

func (s *setup) SetEventConfig(cfg EventConfig) {
	s.eventConfig = cfg
}

func (s *setup) SetCommandConfig(cfg CommandConfig) {
	s.commandConfig = cfg
}

func (s *setup) SetEventBus(bus EventBus) {
	s.eventBus = bus
}

func (s *setup) SetEventStore(store EventStore) {
	s.eventStore = store
}

func (s *setup) SetCommandBus(bus CommandBus) {
	s.commandBus = bus
}

func (s *setup) SetSnapshotRepository(repo SnapshotRepository) {
	s.snapshots = repo
}

func (s *setup) SetAggregateRepository(repo AggregateRepository) {
	s.aggregates = repo
}

func (s *setup) SetEventStoreFactory(f EventStoreFactory) {
	s.eventStoreFactory = f
}

func (s *setup) SetEventBusFactory(f EventBusFactory) {
	s.eventBusFactory = f
}

func (s *setup) SetCommandBusFactory(f CommandBusFactory) {
	s.commandBusFactory = f
}

func (s *setup) AddCommandHandlerFactory(factories ...CommandHandlerFactory) {
	s.commandHandlerFactories = append(s.commandHandlerFactories, factories...)
}

func (s *setup) SetSnapshotRepositoryFactory(f SnapshotRepositoryFactory) {
	s.snapshotRepoFactory = f
}

func (s *setup) SetAggregateRepositoryFactory(f AggregateRepositoryFactory) {
	s.aggregateRepoFactory = f
}

func (s *setup) RegisterAggregate(typ AggregateType, factory AggregateFactory) {
	s.aggregateConfig.Register(typ, factory)
}

func (s *setup) RegisterEvent(typ EventType, factory EventDataFactory) {
	s.eventConfig.Register(typ, factory)
}

func (s *setup) RegisterCommand(typ CommandType, handler CommandHandler) {
	s.commandConfig.Register(typ, handler)
}
