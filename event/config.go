package event

import (
	"encoding/gob"
	"fmt"
	"sync"

	cqrs "github.com/bounoable/cqrs-es"
)

var DisableRegister = false

type config struct {
	mux    sync.RWMutex
	protos map[cqrs.EventType]cqrs.EventData
}

// UnregisteredError is raised when an event type is not registered.
type UnregisteredError struct {
	EventType cqrs.EventType
}

func (err UnregisteredError) Error() string {
	return fmt.Sprintf("unregistered event '%s'", err.EventType)
}

// Config returns a new event config.
func Config() cqrs.EventConfig {
	return &config{
		protos: make(map[cqrs.EventType]cqrs.EventData),
	}
}

func (cfg *config) Register(typ cqrs.EventType, proto cqrs.EventData) {
	if proto == nil {
		panic("nil event data")
	}

	gob.Register(proto)

	cfg.mux.Lock()
	defer cfg.mux.Unlock()

	cfg.protos[typ] = proto
}

func (cfg *config) NewData(typ cqrs.EventType) (cqrs.EventData, error) {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	data, ok := cfg.protos[typ]
	if !ok {
		return nil, UnregisteredError{
			EventType: typ,
		}
	}

	return data, nil
}

func (cfg *config) Protos() map[cqrs.EventType]cqrs.EventData {
	if cfg == nil {
		return map[cqrs.EventType]cqrs.EventData{}
	}

	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[cqrs.EventType]cqrs.EventData)
	for k, v := range cfg.protos {
		m[k] = v
	}

	return m
}
