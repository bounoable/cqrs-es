package channel_test

import (
	"context"
	"testing"

	"github.com/bounoable/cqrs-es/eventbus"
	"github.com/bounoable/cqrs-es/eventbus/channel"
)

func TestPublish(t *testing.T) {
	ctx := context.Background()
	bus := channel.EventBus(ctx)

	eventbus.TestPublish(ctx, t, bus)
}
