package eventbus_test

import (
	"context"
	"testing"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/eventbus"
	"github.com/bounoable/cqrs-es/eventbus/channel"
	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {
	bus := channel.NewEventBus(context.Background(), channel.BufferSize(3))

	ctx := context.Background()
	ch, err := eventbus.Once(ctx, bus, eventbus.TestEventType)
	assert.Nil(t, err)

	data := eventbus.TestEventData{A: 8, B: "test", C: true}
	err = bus.Publish(ctx, cqrs.NewEvent(eventbus.TestEventType, data))
	assert.Nil(t, err)

	evt := <-ch
	assert.Equal(t, data, evt.Data())

	select {
	case <-ch:
	default:
		t.Fatal("channel should be closed")
	}
}
