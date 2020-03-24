package aggregate

//go:generate mockgen -source=config.go -destination=../mocks/aggregate/config.go

import (
	"fmt"
	"sync"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// Factory ...
type Factory func(uuid.UUID) cqrs.Aggregate

// Config ...
type Config interface {
	Register(cqrs.AggregateType, Factory)
	New(cqrs.AggregateType, uuid.UUID) (cqrs.Aggregate, error)
	Factories() map[cqrs.AggregateType]Factory
}

type config struct {
	mux       sync.RWMutex
	factories map[cqrs.AggregateType]Factory
}

// UnregisteredError is raised when an event type is not registered.
type UnregisteredError struct {
	AggregateType cqrs.AggregateType
}

// NewConfig ...
func NewConfig() Config {
	return &config{
		factories: make(map[cqrs.AggregateType]Factory),
	}
}

func (cfg *config) Register(typ cqrs.AggregateType, factory Factory) {
	if factory == nil {
		panic("nil factory")
	}

	cfg.mux.Lock()
	defer cfg.mux.Unlock()
	cfg.factories[typ] = factory
}

func (cfg *config) New(typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error) {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	factory, ok := cfg.factories[typ]
	if !ok {
		return nil, UnregisteredError{
			AggregateType: typ,
		}
	}

	return factory(id), nil
}

func (cfg *config) Factories() map[cqrs.AggregateType]Factory {
	if cfg == nil {
		return map[cqrs.AggregateType]Factory{}
	}

	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[cqrs.AggregateType]Factory)
	for k, v := range cfg.factories {
		m[k] = v
	}

	return m
}

func (err UnregisteredError) Error() string {
	return fmt.Sprintf("unregistered aggregate '%s'", err.AggregateType)
}
