package stan

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/setup"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var (
	defaultConfig = Config{
		BufferSize: 1024,
	}
)

// Config is the events bus config.
type Config struct {
	ClusterID           string
	ClientID            string
	DurableName         string
	URL                 string
	SubjectPrefix       string
	BufferSize          int
	QueueGroup          string
	ConnectOptions      []stan.Option
	SubscriptionOptions []stan.SubscriptionOption
	Logger              *log.Logger
}

// EventBusOption ...
type EventBusOption func(*Config)

type eventBus struct {
	cfg      Config
	eventCfg cqrs.EventConfig
	sc       stan.Conn
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

// ClusterID ...
func ClusterID(id string) EventBusOption {
	return func(cfg *Config) {
		cfg.ClusterID = id
	}
}

// ClientID ...
func ClientID(id string) EventBusOption {
	return func(cfg *Config) {
		cfg.ClientID = id
	}
}

// DurableName ...
func DurableName(name string) EventBusOption {
	return func(cfg *Config) {
		cfg.DurableName = name
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

// BufferSize ...
func BufferSize(size int) EventBusOption {
	return func(cfg *Config) {
		cfg.BufferSize = size
	}
}

// QueueGroup ...
func QueueGroup(group string) EventBusOption {
	return func(cfg *Config) {
		cfg.QueueGroup = group
	}
}

// ConnectOptions ...
func ConnectOptions(options ...stan.Option) EventBusOption {
	return func(cfg *Config) {
		cfg.ConnectOptions = append(cfg.ConnectOptions, options...)
	}
}

// SubscriptionOptions ...
func SubscriptionOptions(options ...stan.SubscriptionOption) EventBusOption {
	return func(cfg *Config) {
		cfg.SubscriptionOptions = append(cfg.SubscriptionOptions, options...)
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
		cfg.URL = stan.DefaultNatsURL
	}

	connectOpts := append([]stan.Option{stan.NatsURL(cfg.URL)}, cfg.ConnectOptions...)
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, connectOpts...)
	if err != nil {
		return nil, err
	}

	return &eventBus{
		cfg:      cfg,
		eventCfg: eventCfg,
		sc:       sc,
	}, nil
}

// NewEventBusWithConnection returns a new NATS streaming event bus.
func NewEventBusWithConnection(sc stan.Conn, eventCfg cqrs.EventConfig, options ...EventBusOption) cqrs.EventBus {
	cfg := defaultConfig

	for _, opt := range options {
		opt(&cfg)
	}

	return &eventBus{
		cfg:      cfg,
		eventCfg: eventCfg,
		sc:       sc,
	}
}

// NewEventBusWithNATSConnection returns a NATS streaming events bus.
func NewEventBusWithNATSConnection(nc *nats.Conn, eventCfg cqrs.EventConfig, options ...EventBusOption) (cqrs.EventBus, error) {
	options = append([]EventBusOption{ConnectOptions(stan.NatsConn(nc))})
	return NewEventBus(eventCfg, options...)
}

// WithEventBusFactory ...
func WithEventBusFactory(options ...EventBusOption) setup.Option {
	return setup.WithEventBusFactory(func(ctx context.Context, s setup.Setup) (cqrs.EventBus, error) {
		return NewEventBus(s.EventConfig(), options...)
	})
}

// WithEventBusFactoryWithConnection ...
func WithEventBusFactoryWithConnection(sc stan.Conn, options ...EventBusOption) setup.Option {
	return setup.WithEventBusFactory(func(ctx context.Context, s setup.Setup) (cqrs.EventBus, error) {
		return NewEventBusWithConnection(sc, s.EventConfig(), options...), nil
	})
}

// WithEventBusFactoryWithNATSConnection ...
func WithEventBusFactoryWithNATSConnection(nc *nats.Conn, options ...EventBusOption) setup.Option {
	return setup.WithEventBusFactory(func(ctx context.Context, s setup.Setup) (cqrs.EventBus, error) {
		return NewEventBusWithNATSConnection(nc, s.EventConfig(), options...)
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

		if err := b.sc.Publish(subject, buf.Bytes()); err != nil {
			return fmt.Errorf("could not publish event: %w", err)
		}
	}

	return nil
}

func (b *eventBus) Subscribe(ctx context.Context, types ...cqrs.EventType) (<-chan cqrs.Event, error) {
	if len(types) == 1 {
		return b.subscribe(ctx, types[0])
	}

	events := make(chan cqrs.Event, b.cfg.BufferSize*len(types))
	var wg sync.WaitGroup

	for _, typ := range types {
		typevents, err := b.subscribe(ctx, typ)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func() {
			for {
				select {
				case <-ctx.Done():
					pushAllEvents(events, typevents)
					wg.Done()
					return
				case event := <-typevents:
					events <- event
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(events)
	}()

	return events, nil
}

func pushAllEvents(target chan<- cqrs.Event, source <-chan cqrs.Event) {
	for event := range source {
		target <- event
	}
}

func (b *eventBus) subscribe(ctx context.Context, typ cqrs.EventType) (<-chan cqrs.Event, error) {
	subject := b.cfg.subject(typ)
	msgs := make(chan *stan.Msg, b.cfg.BufferSize)

	var sub stan.Subscription
	var err error

	options := append([]stan.SubscriptionOption{stan.DurableName(b.cfg.DurableName)}, b.cfg.SubscriptionOptions...)

	if b.cfg.QueueGroup == "" {
		sub, err = b.sc.Subscribe(subject, func(msg *stan.Msg) { msgs <- msg }, options...)
	} else {
		sub, err = b.sc.QueueSubscribe(subject, b.cfg.QueueGroup, func(msg *stan.Msg) { msgs <- msg }, options...)
	}

	if err != nil {
		return nil, fmt.Errorf("stan eventbus: %w", err)
	}

	events := make(chan cqrs.Event, b.cfg.BufferSize)
	handleDone := make(chan struct{})

	go b.handleMessages(msgs, events, handleDone)
	go func() {
		<-ctx.Done()
		if err := sub.Close(); err != nil && b.logger != nil {
			b.logger.Println(fmt.Errorf("stan eventbus: %w", err))
		}

		close(msgs)
		<-handleDone
		close(events)
	}()

	return events, nil
}

func (b *eventBus) handleMessages(msgs <-chan *stan.Msg, events chan<- cqrs.Event, done chan<- struct{}) {
	for msg := range msgs {
		var evtmsg eventMessage
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(&evtmsg); err != nil {
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

		if err := gob.NewDecoder(bytes.NewReader(evtmsg.EventData)).Decode(data); err != nil {
			if b.logger != nil {
				b.logger.Println(err)
			}
			continue
		}

		if evtmsg.AggregateType != cqrs.AggregateType("") && evtmsg.AggregateID != uuid.Nil {
			events <- cqrs.NewAggregateEventWithTime(evtmsg.EventType, data, evtmsg.Time, evtmsg.AggregateType, evtmsg.AggregateID, evtmsg.Version)
		} else {
			events <- cqrs.NewEventWithTime(evtmsg.EventType, data, evtmsg.Time)
		}
	}

	done <- struct{}{}
}

func (cfg Config) subject(typ cqrs.EventType) string {
	return fmt.Sprintf("%s%s", cfg.SubjectPrefix, typ.String())
}
