package cqrs

//go:generate mockgen -source=commandconfig.go -destination=./mocks/commandconfig.go

import (
	"fmt"
	"sync"
)

// CommandConfig ...
type CommandConfig interface {
	Register(CommandType, CommandHandler)
	Handler(CommandType) (CommandHandler, error)
	Handlers() map[CommandType]CommandHandler
}

type commandConfig struct {
	mux      sync.RWMutex
	handlers map[CommandType]CommandHandler
}

// UnregisteredCommandError is raised when a command type is not registered.
type UnregisteredCommandError struct {
	CommandType CommandType
}

func (err UnregisteredCommandError) Error() string {
	return fmt.Sprintf("unregistered command handler '%s'", err.CommandType)
}

// NewCommandConfig returns a new command configuration.
func NewCommandConfig() CommandConfig {
	return &commandConfig{
		handlers: make(map[CommandType]CommandHandler),
	}
}

func (cfg *commandConfig) Register(typ CommandType, handler CommandHandler) {
	cfg.mux.Lock()
	defer cfg.mux.Unlock()
	cfg.handlers[typ] = handler
}

func (cfg *commandConfig) Handler(typ CommandType) (CommandHandler, error) {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	handler, ok := cfg.handlers[typ]
	if !ok {
		return nil, UnregisteredCommandError{
			CommandType: typ,
		}
	}

	return handler, nil
}

func (cfg *commandConfig) Handlers() map[CommandType]CommandHandler {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	m := make(map[CommandType]CommandHandler)
	for k, v := range cfg.handlers {
		m[k] = v
	}

	return m
}
