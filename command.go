package cqrs

//go:generate mockgen -source=command.go -destination=./mocks/command.go

import (
	"context"

	"github.com/google/uuid"
)

// CommandType is a command type.
type CommandType string

// Command is a command.
type Command interface {
	CommandType() CommandType
	AggregateType() AggregateType
	AggregateID() uuid.UUID
}

// CommandBus is the command bus.
type CommandBus interface {
	Dispatch(context.Context, Command) error
}

// CommandConfig ...
type CommandConfig interface {
	Register(CommandType, CommandHandler)
	Handler(CommandType) (CommandHandler, error)
	Handlers() map[CommandType]CommandHandler
}

// CommandHandler handles commands.
type CommandHandler interface {
	HandleCommand(context.Context, Command) error
}
