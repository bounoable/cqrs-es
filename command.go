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
