package command

import (
	"context"

	cqrs "github.com/bounoable/cqrs-es"
)

type bus struct {
	config cqrs.CommandConfig
}

// Bus returns a new CommandBus.
func Bus() cqrs.CommandBus {
	return BusWithConfig(Config())
}

// BusWithConfig ...
func BusWithConfig(config cqrs.CommandConfig) cqrs.CommandBus {
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
