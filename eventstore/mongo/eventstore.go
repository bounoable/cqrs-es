package mongo

import (
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"time"

	"github.com/bounoable/cqrs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// DefaultURI ...
	DefaultURI = "mongodb://localhost:27017"
)

// Config ...
type Config struct {
	URI           string
	Database      string
	Publisher     cqrs.EventPublisher
	ClientOptions []*options.ClientOptions
}

// Option ...
type Option func(*Config)

type eventStore struct {
	config   Config
	eventCfg cqrs.EventConfig
	db       *mongo.Database
}

// Publisher ...
func Publisher(publisher cqrs.EventPublisher) Option {
	return func(cfg *Config) {
		cfg.Publisher = publisher
	}
}

// Database ...
func Database(name string) Option {
	return func(cfg *Config) {
		cfg.Database = name
	}
}

// URI ...
func URI(uri string) Option {
	return func(cfg *Config) {
		cfg.URI = uri
	}
}

// ClientOptions ...
func ClientOptions(options ...*options.ClientOptions) Option {
	return func(cfg *Config) {
		cfg.ClientOptions = append(cfg.ClientOptions, options...)
	}
}

// NewEventStore ...
func NewEventStore(ctx context.Context, eventCfg cqrs.EventConfig, opts ...Option) (cqrs.EventStore, error) {
	var cfg Config
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.URI == "" {
		cfg.URI = os.Getenv("MONGO_EVENTS_URI")
	}

	if cfg.URI == "" {
		cfg.URI = DefaultURI
	}

	if cfg.Database == "" {
		cfg.Database = os.Getenv("MONGO_EVENTS_DB")
	}

	if cfg.Database == "" {
		cfg.Database = "events"
	}

	clientOptions := append([]*options.ClientOptions{options.Client().ApplyURI(cfg.URI)}, cfg.ClientOptions...)
	client, err := mongo.Connect(ctx, clientOptions...)
	if err != nil {
		return nil, wrapError(err)
	}

	return &eventStore{
		db:       client.Database(cfg.Database),
		eventCfg: eventCfg,
	}, nil
}

// WithEventStoreFactory ...
func WithEventStoreFactory(ctx context.Context, options ...Option) cqrs.Option {
	return cqrs.WithEventStoreFactory(func(ctx context.Context, c cqrs.Core) (cqrs.EventStore, error) {
		options = append([]Option{Publisher(c.EventBus())}, options...)
		return NewEventStore(ctx, c.EventConfig(), options...)
	})
}

func (s *eventStore) Save(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, originalVersion int, events ...cqrs.Event) error {
	if len(events) == 0 {
		return nil
	}

	dbEvents := make([]*dbEvent, len(events))

	for i, e := range events {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(e.Data()); err != nil {
			return wrapError(err)
		}

		dbevent := &dbEvent{
			EventType:     e.Type(),
			EventData:     buf.Bytes(),
			Time:          e.Time(),
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			Version:       originalVersion + i + 1,
		}

		dbEvents[i] = dbevent
	}

	docs := make([]interface{}, len(dbEvents))
	for i, dbevent := range dbEvents {
		docs[i] = dbevent
	}

	if err := s.db.Client().UseSession(ctx, func(ctx mongo.SessionContext) error {
		if err := ctx.StartTransaction(); err != nil {
			return err
		}

		res := s.db.Collection("events").FindOne(ctx, bson.M{
			"aggregateType": aggregateType,
			"aggregateId":   aggregateID,
		}, options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}}))

		latestVersion := -1

		var latest dbEvent
		if err := res.Decode(&latest); err == nil {
			latestVersion = latest.Version
		}

		if latestVersion != originalVersion {
			return cqrs.OptimisticConcurrencyError{
				LatestVersion:   latest.Version,
				ProvidedVersion: originalVersion,
			}
		}

		if _, err := s.db.Collection("events").InsertMany(ctx, docs); err != nil {
			return err
		}

		return ctx.CommitTransaction(ctx)
	}); err != nil {
		return wrapError(err)
	}

	if s.config.Publisher != nil {
		if err := s.config.Publisher.Publish(context.Background(), events...); err != nil {
			return wrapError(err)
		}
	}

	return nil
}

func (s *eventStore) Find(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, version int) (cqrs.Event, error) {
	res := s.db.Collection("events").FindOne(ctx, bson.M{
		"aggregateType": aggregateType,
		"aggregateId":   aggregateID,
		"version":       version,
	})

	var dbevent dbEvent
	if err := res.Decode(&dbevent); err != nil {
		return nil, err
	}

	evt, err := s.toCQRSEvent(dbevent)
	if err != nil {
		return nil, wrapError(err)
	}

	return evt, nil
}

func (s *eventStore) Fetch(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, from int, to int) ([]cqrs.Event, error) {
	if from > to {
		return []cqrs.Event{}, nil
	}

	cur, err := s.db.Collection("events").Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
		{Key: "version", Value: bson.D{
			{Key: "$gte", Value: from},
			{Key: "$lte", Value: to},
		}},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, wrapError(err)
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, wrapError(err)
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, wrapError(err)
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) FetchAll(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID) ([]cqrs.Event, error) {
	cur, err := s.db.Collection("events").Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, wrapError(err)
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, wrapError(err)
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, wrapError(err)
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) FetchFrom(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, from int) ([]cqrs.Event, error) {
	cur, err := s.db.Collection("events").Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
		{Key: "version", Value: bson.D{
			{Key: "$gte", Value: from},
		}},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, wrapError(err)
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, wrapError(err)
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, wrapError(err)
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) FetchTo(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, to int) ([]cqrs.Event, error) {
	cur, err := s.db.Collection("events").Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
		{Key: "version", Value: bson.D{
			{Key: "$lte", Value: to},
		}},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, wrapError(err)
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, wrapError(err)
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, wrapError(err)
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) toCQRSEvent(evt dbEvent) (cqrs.Event, error) {
	data, err := s.eventCfg.NewData(evt.EventType)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(evt.EventData)
	if err := gob.NewDecoder(r).Decode(data); err != nil {
		return nil, err
	}

	return cqrs.NewAggregateEvent(evt.EventType, data, evt.Time, evt.AggregateType, evt.AggregateID, evt.Version), nil
}

type dbEvent struct {
	EventType     cqrs.EventType     `bson:"type"`
	EventData     []byte             `bson:"data"`
	Time          time.Time          `bson:"time"`
	AggregateType cqrs.AggregateType `bson:"aggregateType"`
	AggregateID   uuid.UUID          `bson:"aggregateId"`
	Version       int                `bson:"version"`
}

func wrapError(err error) error {
	return wrapError(err)
}
