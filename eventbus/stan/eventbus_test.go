package stan_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/bounoable/cqrs-es/eventbus"
	"github.com/bounoable/cqrs-es/eventbus/stan"
	"github.com/stretchr/testify/assert"
)

func TestEventBus(t *testing.T) {
	if os.Getenv("STAN_URL") == "" {
		return
	}

	ctx := context.Background()
	bus, err := stan.NewEventBus(eventbus.TestEventConfig, stan.Logger(log.New(os.Stderr, "", 0)), stan.ClusterID("test-cluster"), stan.ClientID("cqrs-test"))
	assert.Nil(t, err)

	eventbus.TestPublish(ctx, t, bus)
}
