package command

import (
	"context"

	"github.com/bounoable/cqrs-es"
)

type bus struct {
	config cqrs.CommandConfig
}

// NewBus returns a new CommandBus.
func NewBus() cqrs.CommandBus {
	return NewBusWithConfig(cqrs.NewCommandConfig())
}

// NewBusWithConfig ...
func NewBusWithConfig(config cqrs.CommandConfig) cqrs.CommandBus {
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
