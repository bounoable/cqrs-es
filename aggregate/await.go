package aggregate

import (
	"context"
	"time"

	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// AwaitOption ...
type AwaitOption func(*awaitConfig)

// Await ...
func Await(
	ctx context.Context,
	aggregates cqrs.AggregateRepository,
	typ cqrs.AggregateType,
	id uuid.UUID,
	opts ...AwaitOption,
) (cqrs.Aggregate, error) {
	cfg := newAwaitConfig(opts...)

	cagg := make(chan cqrs.Aggregate, 1)
	cerr := make(chan error, 1)

	var cancel context.CancelFunc
	if cfg.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.timeout)
		defer cancel()
	}

	go func() {
		for {
			agg, err := aggregates.FetchLatest(ctx, typ, id)
			if err != nil {
				cerr <- err
				return
			}

			if len(cfg.matchers) == 0 {
				cagg <- agg
				return
			}

			for _, matcher := range cfg.matchers {
				if matcher(agg) {
					cagg <- agg
					return
				}
			}

			timer := time.NewTimer(cfg.interval)

			select {
			case <-ctx.Done():
				timer.Stop()
				cerr <- ctx.Err()
				return
			case <-timer.C:
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-cerr:
		return nil, err
	case agg := <-cagg:
		return agg, nil
	}
}

// AwaitMatch ...
func AwaitMatch(matchers ...cqrs.AggregateMatcher) AwaitOption {
	return func(cfg *awaitConfig) {
		cfg.matchers = append(cfg.matchers, matchers...)
	}
}

// AwaitInterval ...
func AwaitInterval(dur time.Duration) AwaitOption {
	return func(cfg *awaitConfig) {
		cfg.interval = dur
	}
}

// AwaitTimeout ...
func AwaitTimeout(dur time.Duration) AwaitOption {
	return func(cfg *awaitConfig) {
		cfg.timeout = dur
	}
}

type awaitConfig struct {
	interval time.Duration
	timeout  time.Duration
	matchers []cqrs.AggregateMatcher
}

func newAwaitConfig(opts ...AwaitOption) awaitConfig {
	cfg := awaitConfig{interval: time.Second * 5}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
