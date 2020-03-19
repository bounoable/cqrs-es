package channel

import (
	"context"
	"sync"

	"github.com/bounoable/cqrs-es"
)

// Option ...
type Option func(*config)

type eventBus struct {
	cfg         config
	events      chan cqrs.Event
	mux         sync.RWMutex
	subscribers map[cqrs.EventType][]chan cqrs.Event
}

type config struct {
	BufferSize int
}

// BufferSize ...
func BufferSize(size int) Option {
	return func(cfg *config) {
		cfg.BufferSize = size
	}
}

// NewEventBus ...
func NewEventBus(ctx context.Context, options ...Option) cqrs.EventBus {
	var cfg config
	for _, opt := range options {
		opt(&cfg)
	}

	bus := &eventBus{
		events:      make(chan cqrs.Event, cfg.BufferSize),
		subscribers: make(map[cqrs.EventType][]chan cqrs.Event),
	}

	go bus.start(ctx)

	return bus
}

func (b *eventBus) Publish(ctx context.Context, events ...cqrs.Event) error {
	for _, event := range events {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case b.events <- event:
		}
	}

	return nil
}

func (b *eventBus) Subscribe(ctx context.Context, types ...cqrs.EventType) (<-chan cqrs.Event, error) {
	b.mux.Lock()
	defer b.mux.Unlock()

	var allsubs []chan cqrs.Event

	for _, typ := range types {
		typsubs, ok := b.subscribers[typ]
		if !ok {
			typsubs = make([]chan cqrs.Event, 0)
		}

		sub := make(chan cqrs.Event, b.cfg.BufferSize)
		typsubs = append(typsubs, sub)
		b.subscribers[typ] = typsubs
		allsubs = append(allsubs, typsubs...)
	}

	events := make(chan cqrs.Event, b.cfg.BufferSize*len(allsubs))
	var wg sync.WaitGroup

	for _, sub := range allsubs {
		wg.Add(1)
		go func(sub chan cqrs.Event) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-sub:
					if !ok {
						return
					}
					events <- event
				}
			}
		}(sub)
	}

	go func() {
		wg.Wait()
		close(events)
	}()

	return events, nil
}

func (b *eventBus) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-b.events:
			b.publish(ctx, event)
		}
	}
}

func (b *eventBus) publish(ctx context.Context, event cqrs.Event) {
	b.mux.RLock()
	subs, ok := b.subscribers[event.Type()]
	b.mux.RUnlock()

	if !ok {
		return
	}

	for _, sub := range subs {
		sub <- event
	}
}
