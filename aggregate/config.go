package aggregate

import (
	"fmt"
	"sync"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

type config struct {
	mux       sync.RWMutex
	factories map[cqrs.AggregateType]cqrs.AggregateFactory
}

// UnregisteredError is raised when an event type is not registered.
type UnregisteredError struct {
	AggregateType cqrs.AggregateType
}

// NewConfig ...
func NewConfig() cqrs.AggregateConfig {
	return &config{
		factories: make(map[cqrs.AggregateType]cqrs.AggregateFactory),
	}
}

func (cfg *config) Register(typ cqrs.AggregateType, factory cqrs.AggregateFactory) {
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

func (cfg *config) Factories() map[cqrs.AggregateType]cqrs.AggregateFactory {
	if cfg == nil {
		return map[cqrs.AggregateType]cqrs.AggregateFactory{}
	}

	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[cqrs.AggregateType]cqrs.AggregateFactory)
	for k, v := range cfg.factories {
		m[k] = v
	}

	return m
}

func (err UnregisteredError) Error() string {
	return fmt.Sprintf("unregistered aggregate '%s'", err.AggregateType)
}
