package eventbus

import (
	"context"
	"testing"
	"time"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/event"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	// TestEventType is the EventType used for testing.
	TestEventType = cqrs.EventType("test.event")
	// TestAggregateType is the AggregateType used for testing.
	TestAggregateType = cqrs.AggregateType("test.aggregate")
)

var (
	// TestEventConfig is the EventConfig used for testing.
	TestEventConfig = event.Config()
)

func init() {
	TestEventConfig.Register(TestEventType, TestEventData{})
}

// TestEventData is used for testing EventBus implementations.
type TestEventData struct {
	A int
	B string
	C bool
}

// TestPublish tests the Publish function of an EventBus implementation.
func TestPublish(ctx context.Context, t *testing.T, bus cqrs.EventBus) {
	events, err := bus.Subscribe(ctx, TestEventType)
	assert.Nil(t, err)

	evtch := waitForEvent(TestEventType, events)

	aggregateID := uuid.New()
	data := TestEventData{
		A: 8,
		B: "test",
		C: true,
	}
	err = bus.Publish(ctx, event.NewAggregateEvent(TestEventType, data, TestAggregateType, aggregateID, 6))
	assert.Nil(t, err)

	select {
	case <-time.After(time.Second * 3):
		t.Fatal("event not received")
	case evt := <-evtch:
		assert.Equal(t, data, evt.Data())
	}
}

func waitForEvent(typ cqrs.EventType, events <-chan cqrs.Event) <-chan cqrs.Event {
	ch := make(chan cqrs.Event)
	go func() {
		for evt := range events {
			if evt.Type() == typ {
				ch <- evt
			}
		}
	}()
	return ch
}
