package command

import (
	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// BaseCommand is the base implementation of a command.
type BaseCommand struct {
	typ           cqrs.CommandType
	aggregateType cqrs.AggregateType
	aggregateID   uuid.UUID
}

// NewBase returns a new BaseCommand.
func NewBase(typ cqrs.CommandType, aggregateType cqrs.AggregateType, aggregateID uuid.UUID) BaseCommand {
	return BaseCommand{
		typ:           typ,
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
	}
}

// CommandType ...
func (cmd BaseCommand) CommandType() cqrs.CommandType {
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
