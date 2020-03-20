package eventbus

import (
	"context"
	"time"

	"github.com/bounoable/cqrs-es"
)

// OnceOption ...
type OnceOption func(*onceConfig)

type onceConfig struct {
	timeout  time.Duration
	matchers []Matcher
}

// Matcher ...
type Matcher func(cqrs.Event) bool

// OnceTimeout ...
func OnceTimeout(d time.Duration) OnceOption {
	return func(cfg *onceConfig) {
		cfg.timeout = d
	}
}

// OnceMatch ...
func OnceMatch(matchers ...Matcher) OnceOption {
	return func(cfg *onceConfig) {
		cfg.matchers = matchers
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
		defer cancel()
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case evt, ok := <-events:
				if !ok {
					return
				}

				if len(cfg.matchers) == 0 {
					ch <- evt
					return
				}

				for _, matcher := range cfg.matchers {
					if matcher(evt) {
						ch <- evt
						return
					}
				}
			}
		}
	}()

	return ch, nil
}
