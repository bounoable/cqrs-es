package cqrs

//go:generate mockgen -source=eventconfig.go -destination=./mocks/eventconfig.go

import (
	"encoding/gob"
	"fmt"
	"sync"
)

// EventConfig is the configuration for the events.
type EventConfig interface {
	Register(EventType, EventData)
	// NewData creates an EventData instance for EventType.
	// The returned object is a non-pointer struct.
	NewData(EventType) (EventData, error)
	Protos() map[EventType]EventData
}

type eventConfig struct {
	mux sync.RWMutex
	// factories map[EventType]EventDataFactory
	protos map[EventType]EventData
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
		protos: make(map[EventType]EventData),
	}
}

func (cfg *eventConfig) Register(typ EventType, proto EventData) {
	if proto == nil {
		panic("eventconfig: proto cannot be nil")
	}

	gob.Register(proto)

	cfg.mux.Lock()
	defer cfg.mux.Unlock()

	cfg.protos[typ] = proto
}

func (cfg *eventConfig) NewData(typ EventType) (EventData, error) {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	data, ok := cfg.protos[typ]
	if !ok {
		return nil, UnregisteredEventError{
			EventType: typ,
		}
	}

	return data, nil
}

func (cfg *eventConfig) Protos() map[EventType]EventData {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[EventType]EventData)
	for k, v := range cfg.protos {
		m[k] = v
	}

	return m
}
