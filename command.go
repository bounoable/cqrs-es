package cqrs

//go:generate mockgen -source=command.go -destination=./mocks/command.go

import (
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

// BaseCommand is the base implementation of a command.
type BaseCommand struct {
	typ           CommandType
	aggregateType AggregateType
	aggregateID   uuid.UUID
}

// NewBaseCommand returns a new BaseCommand.
func NewBaseCommand(typ CommandType, aggregateType AggregateType, aggregateID uuid.UUID) BaseCommand {
	return BaseCommand{
		typ:           typ,
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
	}
}

// CommandType ...
func (cmd BaseCommand) CommandType() CommandType {
	return cmd.typ
}

// AggregateType ...
func (cmd BaseCommand) AggregateType() AggregateType {
	return cmd.aggregateType
}

// AggregateID ...
func (cmd BaseCommand) AggregateID() uuid.UUID {
	return cmd.aggregateID
}
