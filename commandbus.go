package cqrs

//go:generate mockgen -source=commandbus.go -destination=./mocks/commandbus.go

import (
	"context"
	"fmt"
	"log"
)

// CommandBus is the command bus.
type CommandBus interface {
	Dispatch(context.Context, Command) error
}

type commandBus struct {
	config CommandConfig
	logger *log.Logger
}

// NewCommandBus returns a new CommandBus.
func NewCommandBus(logger *log.Logger) CommandBus {
	return NewCommandBusWithConfig(NewCommandConfig(), logger)
}

// NewCommandBusWithConfig ...
func NewCommandBusWithConfig(config CommandConfig, logger *log.Logger) CommandBus {
	return &commandBus{
		config: config,
		logger: logger,
	}
}

func (b *commandBus) Dispatch(ctx context.Context, cmd Command) error {
	handler, err := b.config.Handler(cmd.CommandType())
	if err != nil {
		if b.logger != nil {
			b.logger.Println(fmt.Sprintf("commandbus (Dispatch): %s", err))
		}

		return err
	}

	if err := handler.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	return nil
}
