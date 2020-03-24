package command

//go:generate mockgen -source=config.go -destination=../mocks/command/config.go

import (
	"context"
	"fmt"
	"sync"
)

// Config ...
type Config interface {
	Register(Type, Handler)
	Handler(Type) (Handler, error)
	Handlers() map[Type]Handler
}

// Handler handles commands.
type Handler interface {
	HandleCommand(context.Context, Command) error
}

type config struct {
	mux      sync.RWMutex
	handlers map[Type]Handler
}

// UnregisteredError is raised when a command type is not registered.
type UnregisteredError struct {
	CommandType Type
}

func (err UnregisteredError) Error() string {
	return fmt.Sprintf("unregistered command handler '%s'", err.CommandType)
}

// NewConfig returns a new command configuration.
func NewConfig() Config {
	return &config{
		handlers: make(map[Type]Handler),
	}
}

func (cfg *config) Register(typ Type, handler Handler) {
	cfg.mux.Lock()
	defer cfg.mux.Unlock()
	cfg.handlers[typ] = handler
}

func (cfg *config) Handler(typ Type) (Handler, error) {
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

func (cfg *config) Handlers() map[Type]Handler {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[Type]Handler)
	for k, v := range cfg.handlers {
		m[k] = v
	}

	return m
}
