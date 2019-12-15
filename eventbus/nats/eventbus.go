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

var (
	defaultConfig = Config{
		SubscriptionBufferSize: 1024,
	}
)

// Config is the events bus config.
type Config struct {
	URL                    string
	SubjectPrefix          string
	SubscriptionBufferSize int
	QueueGroup             string
	ConnectOptions         []nats.Option
	Logger                 *log.Logger
}

// EventBusOption ...
type EventBusOption func(*Config)

type eventBus struct {
	cfg      Config
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

// Logger ...
func Logger(logger *log.Logger) EventBusOption {
	return func(cfg *Config) {
		cfg.Logger = logger
	}
}

// URL ...
func URL(url string) EventBusOption {
	return func(cfg *Config) {
		cfg.URL = url
	}
}

// SubjectPrefix ...
func SubjectPrefix(prefix string) EventBusOption {
	return func(cfg *Config) {
		cfg.SubjectPrefix = prefix
	}
}

// SubscriptionBufferSize ...
func SubscriptionBufferSize(size int) EventBusOption {
	return func(cfg *Config) {
		cfg.SubscriptionBufferSize = size
	}
}

// QueueGroup ...
func QueueGroup(group string) EventBusOption {
	return func(cfg *Config) {
		cfg.QueueGroup = group
	}
}

// ConnectOptions ...
func ConnectOptions(options ...nats.Option) EventBusOption {
	return func(cfg *Config) {
		cfg.ConnectOptions = append(cfg.ConnectOptions, options...)
	}
}

// NewEventBus ...
func NewEventBus(eventCfg cqrs.EventConfig, options ...EventBusOption) (cqrs.EventBus, error) {
	cfg := defaultConfig

	for _, opt := range options {
		opt(&cfg)
	}

	if cfg.URL == "" {
		cfg.URL = os.Getenv("NATS_URL")
	}

	if cfg.URL == "" {
		cfg.URL = nats.DefaultURL
	}

	nc, err := nats.Connect(cfg.URL, cfg.ConnectOptions...)
	if err != nil {
		return nil, err
	}

	return &eventBus{
		cfg:      cfg,
		eventCfg: eventCfg,
		nc:       nc,
	}, nil
}

// NewEventBusWithConnection returns a new NATS event bus.
func NewEventBusWithConnection(nc *nats.Conn, eventCfg cqrs.EventConfig, options ...EventBusOption) cqrs.EventBus {
	cfg := defaultConfig

	for _, opt := range options {
		opt(&cfg)
	}

	return &eventBus{
		cfg:      cfg,
		eventCfg: eventCfg,
		nc:       nc,
	}
}

// WithEventBusFactory ...
func WithEventBusFactory(options ...EventBusOption) cqrs.Option {
	return cqrs.WithEventBusFactory(func(ctx context.Context, c cqrs.Core) (cqrs.EventBus, error) {
		return NewEventBus(c.EventConfig(), options...)
	})
}

// WithEventBusFactoryWithConnection ...
func WithEventBusFactoryWithConnection(nc *nats.Conn, options ...EventBusOption) cqrs.Option {
	return cqrs.WithEventBusFactory(func(ctx context.Context, c cqrs.Core) (cqrs.EventBus, error) {
		return NewEventBusWithConnection(nc, c.EventConfig(), options...), nil
	})
}

func (b *eventBus) Publish(_ context.Context, events ...cqrs.Event) error {
	for _, e := range events {
		var dataBuf bytes.Buffer
		if err := gob.NewEncoder(&dataBuf).Encode(e.Data()); err != nil {
			return err
		}

		subject := b.cfg.subject(e.Type())
		evt := &eventMessage{
			EventType:     e.Type(),
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

	var sub *nats.Subscription
	var err error

	if b.cfg.QueueGroup == "" {
		sub, err = b.nc.ChanSubscribe(subject, msgs)
	} else {
		sub, err = b.nc.ChanQueueSubscribe(subject, b.cfg.QueueGroup, msgs)
	}

	if err != nil {
		return nil, fmt.Errorf("nats eventbus: %w", err)
	}

	events := make(chan cqrs.Event, b.cfg.SubscriptionBufferSize)
	go b.handleMessages(msgs, events)
	go func() {
		<-ctx.Done()
		if err := sub.Unsubscribe(); err != nil && b.logger != nil {
			b.logger.Println(fmt.Errorf("nats eventbus: %w", err))
		}

		close(msgs)
		close(events)
	}()

	return events, nil
}

func (b *eventBus) handleMessages(msgs <-chan *nats.Msg, events chan<- cqrs.Event) {
	for msg := range msgs {
		var evtmsg eventMessage
		if err := gob.NewDecoder(bytes.NewBuffer(msg.Data)).Decode(&evtmsg); err != nil {
			if b.logger != nil {
				b.logger.Println(err)
			}
			continue
		}

		data, err := b.eventCfg.NewData(evtmsg.EventType)
		if err != nil {
			if b.logger != nil {
				b.logger.Println(err)
			}
			continue
		}

		if err := gob.NewDecoder(bytes.NewBuffer(evtmsg.EventData)).Decode(data); err != nil {
			if b.logger != nil {
				b.logger.Println(err)
			}
			continue
		}

		if evtmsg.AggregateType != cqrs.AggregateType("") && evtmsg.AggregateID != uuid.Nil {
			events <- cqrs.NewAggregateEvent(evtmsg.EventType, data, evtmsg.Time, evtmsg.AggregateType, evtmsg.AggregateID, evtmsg.Version)
		} else {
			events <- cqrs.NewEvent(evtmsg.EventType, data, evtmsg.Time)
		}
	}
}

func (cfg Config) subject(typ cqrs.EventType) string {
	return fmt.Sprintf("%s%s", cfg.SubjectPrefix, typ.String())
}
