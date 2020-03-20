package eventbus

import (
	"context"
	"time"

	"github.com/bounoable/cqrs-es"
)

// OnceOption ...
type OnceOption func(*onceConfig)

type onceConfig struct {
	timeout time.Duration
}

// OnceTimeout ...
func OnceTimeout(d time.Duration) OnceOption {
	return func(cfg *onceConfig) {
		cfg.timeout = d
	}
}

// Once ...
func Once(ctx context.Context, bus cqrs.EventBus, typ cqrs.EventType, opts ...OnceOption) (<-chan cqrs.Event, error) {
	var cfg onceConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var cancel context.CancelFunc
	if cfg.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	events, err := bus.Subscribe(ctx, typ)
	if err != nil {
		cancel()
		return nil, err
	}

	ch := make(chan cqrs.Event, 1)
	go func() {
		defer close(ch)
		defer cancel()
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-events:
			if !ok {
				return
			}
			ch <- evt
		}
	}()

	return ch, nil
}
