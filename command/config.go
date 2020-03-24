package command

//go:generate mockgen -source=config.go -destination=../mocks/command/config.go

import (
	"context"
	"fmt"
	"sync"

	cqrs "github.com/bounoable/cqrs-es"
)

// Config ...
type Config interface {
	Register(cqrs.CommandType, Handler)
	Handler(cqrs.CommandType) (Handler, error)
	Handlers() map[cqrs.CommandType]Handler
}

// Handler handles commands.
type Handler interface {
	HandleCommand(context.Context, cqrs.Command) error
}

type config struct {
	mux      sync.RWMutex
	handlers map[cqrs.CommandType]Handler
}

// UnregisteredError is raised when a command type is not registered.
type UnregisteredError struct {
	CommandType cqrs.CommandType
}

func (err UnregisteredError) Error() string {
	return fmt.Sprintf("unregistered command handler '%s'", err.CommandType)
}

// NewConfig returns a new command configuration.
func NewConfig() Config {
	return &config{
		handlers: make(map[cqrs.CommandType]Handler),
	}
}

func (cfg *config) Register(typ cqrs.CommandType, handler Handler) {
	cfg.mux.Lock()
	defer cfg.mux.Unlock()
	cfg.handlers[typ] = handler
}

func (cfg *config) Handler(typ cqrs.CommandType) (Handler, error) {
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

func (cfg *config) Handlers() map[cqrs.CommandType]Handler {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[cqrs.CommandType]Handler)
	for k, v := range cfg.handlers {
		m[k] = v
	}

	return m
}
