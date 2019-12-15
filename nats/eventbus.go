package nats

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bounoable/cqrs"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// EventBusConfig is the events bus config.
type EventBusConfig struct {
	SubjectPrefix          string
	SubscriptionBufferSize int
}

// EventBusOption ...
type EventBusOption func(*eventBus)

type eventBus struct {
	cfg      EventBusConfig
	eventCfg cqrs.EventConfig
	nc       *nats.Conn
	logger   *log.Logger
}

type eventMessage struct {
	EventType     cqrs.EventType
	EventData     []byte
	Time          time.Time
	AggregateType cqrs.AggregateType
	AggregateID   uuid.UUID
	Version       int
}

// WithLogger ...
func WithLogger(logger *log.Logger) EventBusOption {
	return func(bus *eventBus) {
		bus.logger = logger
	}
}

// NewEventBus ...
func NewEventBus(cfg EventBusConfig, eventCfg cqrs.EventConfig, options ...EventBusOption) (cqrs.EventBus, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return NewEventBusWithConnection(nc, cfg, eventCfg, options...), nil
}

// NewEventBusWithConnection returns a new NATS event bus.
func NewEventBusWithConnection(nc *nats.Conn, cfg EventBusConfig, eventCfg cqrs.EventConfig, options ...EventBusOption) cqrs.EventBus {
	bus := &eventBus{
		cfg:      cfg,
		eventCfg: eventCfg,
		nc:       nc,
	}

	for _, opt := range options {
		opt(bus)
	}

	return bus
}

// WithEventBusFactory ...
func WithEventBusFactory(cfg EventBusConfig) cqrs.Option {
	return cqrs.WithEventBusFactory(func(ctx context.Context, c cqrs.Core) (cqrs.EventBus, error) {
		return NewEventBus(cfg, c.EventConfig())
	})
}

// WithEventBusFactoryWithConnection ...
func WithEventBusFactoryWithConnection(cfg EventBusConfig, nc *nats.Conn) cqrs.Option {
	return cqrs.WithEventBusFactory(func(ctx context.Context, c cqrs.Core) (cqrs.EventBus, error) {
		return NewEventBusWithConnection(nc, cfg, c.EventConfig()), nil
	})
}

func (b *eventBus) Publish(ctx context.Context, events ...cqrs.Event) error {
	for _, e := range events {
		var dataBuf bytes.Buffer
		if err := gob.NewEncoder(&dataBuf).Encode(e.EventData()); err != nil {
			return err
		}

		subject := b.cfg.subject(e.EventType())
		evt := &eventMessage{
			EventType:     e.EventType(),
			EventData:     dataBuf.Bytes(),
			Time:          e.Time(),
			AggregateType: e.AggregateType(),
			AggregateID:   e.AggregateID(),
			Version:       e.Version(),
		}

		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(evt); err != nil {
			return fmt.Errorf("could not marshal event data: %w", err)
		}

		if err := b.nc.Publish(subject, buf.Bytes()); err != nil {
			return fmt.Errorf("could not publish event: %w", err)
		}
	}

	return nil
}

func (b *eventBus) Subscribe(ctx context.Context, typ cqrs.EventType) (<-chan cqrs.Event, error) {
	subject := b.cfg.subject(typ)
	msgs := make(chan *nats.Msg, b.cfg.SubscriptionBufferSize)

	_, err := b.nc.ChanSubscribe(subject, msgs)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe: %w", err)
	}

	events := make(chan cqrs.Event, b.cfg.SubscriptionBufferSize)
	go b.handleMessages(ctx, msgs, events)

	return events, nil
}

func (b *eventBus) handleMessages(ctx context.Context, msgs <-chan *nats.Msg, events chan<- cqrs.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			var evtmsg eventMessage
			if err := gob.NewDecoder(bytes.NewBuffer(msg.Data)).Decode(&evtmsg); err != nil {
				if b.logger != nil {
					b.logger.Println(err)
				}
				break
			}

			data, err := b.eventCfg.NewData(evtmsg.EventType)
			if err != nil {
				if b.logger != nil {
					b.logger.Println(err)
				}
				break
			}

			if err := gob.NewDecoder(bytes.NewBuffer(evtmsg.EventData)).Decode(data); err != nil {
				if b.logger != nil {
					b.logger.Println(err)
				}
				break
			}

			if evtmsg.AggregateType != cqrs.AggregateType("") && evtmsg.AggregateID != uuid.Nil {
				events <- cqrs.NewAggregateEvent(evtmsg.EventType, data, evtmsg.Time, evtmsg.AggregateType, evtmsg.AggregateID, evtmsg.Version)
			} else {
				events <- cqrs.NewEvent(evtmsg.EventType, data, evtmsg.Time)
			}
		}
	}
}

func (cfg EventBusConfig) subject(typ cqrs.EventType) string {
	return fmt.Sprintf("%s%s", cfg.SubjectPrefix, typ.String())
}
