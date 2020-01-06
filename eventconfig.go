package cqrs

//go:generate mockgen -source=eventconfig.go -destination=./mocks/eventconfig.go

import (
	"fmt"
	"reflect"
	"sync"
)

// EventDataFactory ...
type EventDataFactory func() EventData

// EventConfig is the configuration for the events.
type EventConfig interface {
	Register(EventType, EventData)
	// NewData creates an EventData instance for EventType.
	// The returned object is a struct.
	NewData(EventType) (EventData, error)
	Factories() map[EventType]EventDataFactory
}

type eventConfig struct {
	mux       sync.RWMutex
	factories map[EventType]EventDataFactory
}

// UnregisteredEventError is raised when an event type is not registered.
type UnregisteredEventError struct {
	EventType EventType
}

func (err UnregisteredEventError) Error() string {
	return fmt.Sprintf("unregistered event '%s'", err.EventType)
}

// NewEventConfig returns a new event config.
func NewEventConfig() EventConfig {
	return &eventConfig{
		factories: make(map[EventType]EventDataFactory),
	}
}

func (cfg *eventConfig) Register(typ EventType, proto EventData) {
	if proto == nil {
		panic("eventconfig: proto cannot be nil")
	}

	refval := reflect.TypeOf(proto)
	for refval.Kind() == reflect.Ptr {
		refval = refval.Elem()
	}

	cfg.RegisterFactory(typ, func() EventData {
		return reflect.New(refval).Interface()
	})
}

func (cfg *eventConfig) RegisterFactory(typ EventType, factory EventDataFactory) {
	if factory == nil {
		panic("eventconfig: factory cannot be nil")
	}

	cfg.mux.Lock()
	defer cfg.mux.Unlock()
	cfg.factories[typ] = factory
}

func (cfg *eventConfig) NewData(typ EventType) (EventData, error) {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	factory, ok := cfg.factories[typ]
	if !ok {
		return nil, UnregisteredEventError{
			EventType: typ,
		}
	}

	return factory(), nil
}

func (cfg *eventConfig) Factories() map[EventType]EventDataFactory {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[EventType]EventDataFactory)
	for k, v := range cfg.factories {
		m[k] = v
	}

	return m
}
