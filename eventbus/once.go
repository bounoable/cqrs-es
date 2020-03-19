package eventbus

import (
	"context"

	"github.com/bounoable/cqrs-es"
)

// Once ...
func Once(ctx context.Context, bus cqrs.EventBus, typ cqrs.EventType) (<-chan cqrs.Event, error) {
	ch := make(chan cqrs.Event, 1)

	ctx, cancel := context.WithCancel(ctx)

	events, err := bus.Subscribe(ctx, typ)
	if err != nil {
		cancel()
		return nil, err
	}

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
