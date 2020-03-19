package nats_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/bounoable/cqrs-es/eventbus"
	"github.com/bounoable/cqrs-es/eventbus/nats"
	"github.com/stretchr/testify/assert"
)

func TestEventBus(t *testing.T) {
	if os.Getenv("NATS_URL") == "" {
		return
	}

	ctx := context.Background()
	bus, err := nats.NewEventBus(eventbus.TestEventConfig, nats.Logger(log.New(os.Stderr, "", 0)))
	assert.Nil(t, err)

	eventbus.TestPublish(ctx, t, bus)
}
