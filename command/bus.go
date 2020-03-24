package command

import (
	"context"

	"github.com/bounoable/cqrs-es"
)

type bus struct {
	config Config
}

// NewBus returns a new CommandBus.
func NewBus() cqrs.CommandBus {
	return NewBusWithConfig(NewConfig())
}

// NewBusWithConfig ...
func NewBusWithConfig(config Config) cqrs.CommandBus {
	return &bus{
		config: config,
	}
}

func (b *bus) Dispatch(ctx context.Context, cmd cqrs.Command) error {
	handler, err := b.config.Handler(cmd.CommandType())
	if err != nil {
		return err
	}

	if err := handler.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	return nil
}
