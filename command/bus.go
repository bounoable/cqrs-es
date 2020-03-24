package command

//go:generate mockgen -source=bus.go -destination=../mocks/command/bus.go

import (
	"context"
)

// Bus is the command bus.
type Bus interface {
	Dispatch(context.Context, Command) error
}

type bus struct {
	config Config
}

// NewBus returns a new CommandBus.
func NewBus() Bus {
	return NewBusWithConfig(NewConfig())
}

// NewBusWithConfig ...
func NewBusWithConfig(config Config) Bus {
	return &bus{
		config: config,
	}
}

func (b *bus) Dispatch(ctx context.Context, cmd Command) error {
	handler, err := b.config.Handler(cmd.CommandType())
	if err != nil {
		return err
	}

	if err := handler.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	return nil
}
