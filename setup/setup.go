package setup

//go:generate mockgen -source=setup.go -destination=../mocks/setup/setup.go

import (
	"context"
	"errors"
	"fmt"
	"log"

	cqrs "github.com/bounoable/cqrs-es"
)

var (
	globalAggregateConfig cqrs.AggregateConfig
	globalEventConfig     cqrs.EventConfig
)

// EventStoreFactory ...
type EventStoreFactory func(context.Context, Setup) (cqrs.EventStore, error)

// EventBusFactory ...
type EventBusFactory func(context.Context, Setup) (cqrs.EventBus, error)

// CommandBusFactory ...
type CommandBusFactory func(context.Context, Setup) (cqrs.CommandBus, error)

// CommandHandlerFactory ...
type CommandHandlerFactory func(context.Context, Setup) (cqrs.CommandHandler, []cqrs.CommandType, error)

// SnapshotRepositoryFactory ...
type SnapshotRepositoryFactory func(context.Context, Setup) (cqrs.SnapshotRepository, error)

// AggregateRepositoryFactory ...
type AggregateRepositoryFactory func(context.Context, Setup) (cqrs.AggregateRepository, error)

// Option ...
type Option func(Setup)

// Setup ...
type Setup interface {
	SetLogger(*log.Logger)

	AggregateConfig() cqrs.AggregateConfig
	EventConfig() cqrs.EventConfig
	CommandConfig() cqrs.CommandConfig

	SetEventStoreFactory(EventStoreFactory)
	SetEventBusFactory(EventBusFactory)
	SetCommandBusFactory(CommandBusFactory)
	AddCommandHandlerFactory(...CommandHandlerFactory)
	SetSnapshotRepositoryFactory(SnapshotRepositoryFactory)
	SetAggregateRepositoryFactory(AggregateRepositoryFactory)

	EventStore() cqrs.EventStore
	EventBus() cqrs.EventBus
	CommandBus() cqrs.CommandBus
	SnapshotRepository() cqrs.SnapshotRepository
	AggregateRepository() cqrs.AggregateRepository

	RegisterAggregate(cqrs.AggregateType, cqrs.AggregateFactory)
	RegisterEvent(cqrs.EventType, cqrs.EventData)
	RegisterCommand(cqrs.CommandType, cqrs.CommandHandler)
}

// RegisterAggregate ...
func RegisterAggregate(typ cqrs.AggregateType, factory cqrs.AggregateFactory) {
	globalAggregateConfig.Register(typ, factory)
}

// RegisterEvent ...
func RegisterEvent(typ cqrs.EventType, proto cqrs.EventData) {
	globalEventConfig.Register(typ, proto)
}

// New ...
func New(ctx context.Context, options ...Option) (cqrs.Container, error) {
	return NewWithConfigs(ctx, cqrs.NewAggregateConfig(), cqrs.NewEventConfig(), cqrs.NewCommandConfig(), options...)
}

// NewWithConfigs ...
func NewWithConfigs(
	ctx context.Context,
	aggregateConfig cqrs.AggregateConfig,
	eventConfig cqrs.EventConfig,
	commandConfig cqrs.CommandConfig,
	options ...Option,
) (cqrs.Container, error) {
	if aggregateConfig == nil {
		return nil, errors.New("aggregate config cannot be nil")
	}

	if eventConfig == nil {
		return nil, errors.New("event config cannot be nil")
	}

	if commandConfig == nil {
		return nil, errors.New("command config cannot be nil")
	}

	s := &setup{
		aggregateConfig: aggregateConfig,
		eventConfig:     eventConfig,
		commandConfig:   commandConfig,
	}

	applyGlobalConfigs(aggregateConfig, eventConfig)

	for _, opt := range options {
		opt(s)
	}

	if s.eventBusFactory != nil {
		bus, err := s.eventBusFactory(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("setup: %w", err)
		}
		s.eventBus = bus
	} else if s.logger != nil {
		s.logger.Println("setup: no event bus")
	}

	if s.eventStoreFactory != nil {
		store, err := s.eventStoreFactory(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("setup: %w", err)
		}
		s.eventStore = store
	} else if s.logger != nil {
		s.logger.Println("setup: no event store")
	}

	if s.commandBusFactory != nil {
		bus, err := s.commandBusFactory(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("setup: %w", err)
		}
		s.commandBus = bus
	} else if s.logger != nil {
		s.logger.Println("setup: no command bus")
	}

	if s.snapshotRepoFactory != nil {
		repo, err := s.snapshotRepoFactory(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("setup: %w", err)
		}
		s.snapshotRepo = repo
	} else if s.logger != nil {
		s.logger.Println("setup: no snapshot repository")
	}

	if s.aggregateRepoFactory != nil {
		repo, err := s.aggregateRepoFactory(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("setup: %w", err)
		}
		s.aggregateRepo = repo
	} else if s.eventStore != nil && s.aggregateConfig != nil {
		s.aggregateRepo = cqrs.NewAggregateRepository(s.eventStore, s.aggregateConfig, s.snapshotRepo)
	} else if s.logger != nil {
		s.logger.Println("setup: no aggregate repository")
	}

	for _, f := range s.commandHandlerFactories {
		handler, types, err := f(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("setup: %w", err)
		}

		for _, typ := range types {
			s.commandConfig.Register(typ, handler)
		}
	}

	return cqrs.New(
		s.logger,
		s.aggregateConfig,
		s.eventConfig,
		s.commandConfig,
		s.eventBus,
		s.eventStore,
		s.commandBus,
		s.snapshotRepo,
		s.aggregateRepo,
	), nil
}

// WithLogger ...
func WithLogger(logger *log.Logger) Option {
	return func(s Setup) {
		s.SetLogger(logger)
	}
}

// WithAggregate ...
func WithAggregate(typ cqrs.AggregateType, factory cqrs.AggregateFactory) Option {
	return func(s Setup) {
		s.RegisterAggregate(typ, factory)
	}
}

// WithEvent ...
func WithEvent(typ cqrs.EventType, proto cqrs.EventData) Option {
	return func(s Setup) {
		s.RegisterEvent(typ, proto)
	}
}

// WithCommand ...
func WithCommand(typ cqrs.CommandType, handler cqrs.CommandHandler) Option {
	return func(s Setup) {
		s.RegisterCommand(typ, handler)
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

type setup struct {
	logger *log.Logger

	aggregateConfig cqrs.AggregateConfig
	eventConfig     cqrs.EventConfig
	commandConfig   cqrs.CommandConfig

	eventStoreFactory       EventStoreFactory
	eventBusFactory         EventBusFactory
	commandBusFactory       CommandBusFactory
	commandHandlerFactories []CommandHandlerFactory
	snapshotRepoFactory     SnapshotRepositoryFactory
	aggregateRepoFactory    AggregateRepositoryFactory

	eventStore    cqrs.EventStore
	eventBus      cqrs.EventBus
	commandBus    cqrs.CommandBus
	snapshotRepo  cqrs.SnapshotRepository
	aggregateRepo cqrs.AggregateRepository
}

func (s *setup) AggregateConfig() cqrs.AggregateConfig {
	return s.aggregateConfig
}

func (s *setup) EventConfig() cqrs.EventConfig {
	return s.eventConfig
}

func (s *setup) CommandConfig() cqrs.CommandConfig {
	return s.commandConfig
}

func (s *setup) SetLogger(l *log.Logger) {
	s.logger = l
}

func (s *setup) SetAggregateConfig(cfg cqrs.AggregateConfig) {
	s.aggregateConfig = cfg
}

func (s *setup) SetEventConfig(cfg cqrs.EventConfig) {
	s.eventConfig = cfg
}

func (s *setup) SetCommandConfig(cfg cqrs.CommandConfig) {
	s.commandConfig = cfg
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

func (s *setup) EventStore() cqrs.EventStore {
	return s.eventStore
}

func (s *setup) EventBus() cqrs.EventBus {
	return s.eventBus
}

func (s *setup) CommandBus() cqrs.CommandBus {
	return s.commandBus
}

func (s *setup) SnapshotRepository() cqrs.SnapshotRepository {
	return s.snapshotRepo
}

func (s *setup) AggregateRepository() cqrs.AggregateRepository {
	return s.aggregateRepo
}

func (s *setup) RegisterAggregate(typ cqrs.AggregateType, factory cqrs.AggregateFactory) {
	s.aggregateConfig.Register(typ, factory)
}

func (s *setup) RegisterEvent(typ cqrs.EventType, proto cqrs.EventData) {
	s.eventConfig.Register(typ, proto)
}

func (s *setup) RegisterCommand(typ cqrs.CommandType, handler cqrs.CommandHandler) {
	s.commandConfig.Register(typ, handler)
}

func applyGlobalConfigs(aggregates cqrs.AggregateConfig, events cqrs.EventConfig) {
	for k, v := range globalAggregateConfig.Factories() {
		aggregates.Register(k, v)
	}

	for k, v := range globalEventConfig.Factories() {
		events.Register(k, v())
	}
}
