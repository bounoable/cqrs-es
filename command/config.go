package command

import (
	"fmt"
	"sync"

	cqrs "github.com/bounoable/cqrs-es"
)

type config struct {
	mux      sync.RWMutex
	handlers map[cqrs.CommandType]cqrs.CommandHandler
}

// UnregisteredError is raised when a command type is not registered.
type UnregisteredError struct {
	CommandType cqrs.CommandType
}

func (err UnregisteredError) Error() string {
	return fmt.Sprintf("unregistered command handler '%s'", err.CommandType)
}

// Config returns a new command configuration.
func Config() cqrs.CommandConfig {
	return &config{
		handlers: make(map[cqrs.CommandType]cqrs.CommandHandler),
	}
}

func (cfg *config) Register(typ cqrs.CommandType, handler cqrs.CommandHandler) {
	cfg.mux.Lock()
	defer cfg.mux.Unlock()
	cfg.handlers[typ] = handler
}

func (cfg *config) Handler(typ cqrs.CommandType) (cqrs.CommandHandler, error) {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	handler, ok := cfg.handlers[typ]
	if !ok {
		return nil, UnregisteredError{
			CommandType: typ,
		}
	}

	return handler, nil
}

func (cfg *config) Handlers() map[cqrs.CommandType]cqrs.CommandHandler {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[cqrs.CommandType]cqrs.CommandHandler)
	for k, v := range cfg.handlers {
		m[k] = v
	}

	return m
}
