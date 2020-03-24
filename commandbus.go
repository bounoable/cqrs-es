package cqrs

//go:generate mockgen -source=commandbus.go -destination=./mocks/commandbus.go

import (
	"context"
)

// CommandBus is the command bus.
type CommandBus interface {
	Dispatch(context.Context, Command) error
}
