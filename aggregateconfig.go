package cqrs

//go:generate mockgen -source=aggregateconfig.go -destination=./mocks/aggregateconfig.go

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// AggregateFactory ...
type AggregateFactory func(uuid.UUID) Aggregate

// AggregateConfig ...
type AggregateConfig interface {
	Register(AggregateType, AggregateFactory)
	New(AggregateType, uuid.UUID) (Aggregate, error)
	Factories() map[AggregateType]AggregateFactory
}

type aggregateConfig struct {
	mux       sync.RWMutex
	factories map[AggregateType]AggregateFactory
}

// UnregisteredAggregateError is raised when an event type is not registered.
type UnregisteredAggregateError struct {
	AggregateType AggregateType
}

// NewAggregateConfig ...
func NewAggregateConfig() AggregateConfig {
	return &aggregateConfig{
		factories: make(map[AggregateType]AggregateFactory),
	}
}

func (cfg *aggregateConfig) Register(typ AggregateType, factory AggregateFactory) {
	if factory == nil {
		panic("aggregate config: aggregate factory cannot be nil")
	}

	cfg.mux.Lock()
	defer cfg.mux.Unlock()
	cfg.factories[typ] = factory
}

func (cfg *aggregateConfig) New(typ AggregateType, id uuid.UUID) (Aggregate, error) {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	factory, ok := cfg.factories[typ]
	if !ok {
		return nil, UnregisteredAggregateError{
			AggregateType: typ,
		}
	}

	return factory(id), nil
}

func (cfg *aggregateConfig) Factories() map[AggregateType]AggregateFactory {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[AggregateType]AggregateFactory)
	for k, v := range cfg.factories {
		m[k] = v
	}

	return m
}

func (err UnregisteredAggregateError) Error() string {
	return fmt.Sprintf("unregistered aggregate '%s'", err.AggregateType)
}
