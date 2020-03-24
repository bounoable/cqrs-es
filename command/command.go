package command

//go:generate mockgen -source=command.go -destination=../mocks/command/command.go

import (
	cqrs "github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// Type is a command type.
type Type string

// Command is a command.
type Command interface {
	CommandType() Type
	AggregateType() cqrs.AggregateType
	AggregateID() uuid.UUID
}

// BaseCommand is the base implementation of a command.
type BaseCommand struct {
	typ           Type
	aggregateType cqrs.AggregateType
	aggregateID   uuid.UUID
}

// NewBase returns a new BaseCommand.
func NewBase(typ Type, aggregateType cqrs.AggregateType, aggregateID uuid.UUID) BaseCommand {
	return BaseCommand{
		typ:           typ,
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
	}
}

// CommandType ...
func (cmd BaseCommand) CommandType() Type {
	return cmd.typ
}

// AggregateType ...
func (cmd BaseCommand) AggregateType() cqrs.AggregateType {
	return cmd.aggregateType
}

// AggregateID ...
func (cmd BaseCommand) AggregateID() uuid.UUID {
	return cmd.aggregateID
}
