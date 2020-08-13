package mongo

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/event"
	"github.com/bounoable/cqrs-es/setup"
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
	URI              string
	Database         string
	Publisher        cqrs.EventPublisher
	ResolvePublisher func() cqrs.EventPublisher
	Transactions     bool
	CreateIndexes    bool
	ClientOptions    []*options.ClientOptions
}

// Option ...
type Option func(*Config)

type eventStore struct {
	config    Config
	eventCfg  cqrs.EventConfig
	db        *mongo.Database
	col       *mongo.Collection
	publisher cqrs.EventPublisher
}

// Publisher ...
func Publisher(publisher cqrs.EventPublisher) Option {
	return func(cfg *Config) {
		cfg.Publisher = publisher
	}
}

// ResolvePublisher ...
func ResolvePublisher(resolve func() cqrs.EventPublisher) Option {
	return func(cfg *Config) {
		cfg.ResolvePublisher = resolve
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

// Transactions ...
func Transactions(use bool) Option {
	return func(cfg *Config) {
		cfg.Transactions = use
	}
}

// CreateIndexes ...
func CreateIndexes() Option {
	return func(cfg *Config) {
		cfg.CreateIndexes = true
	}
}

// ClientOptions ...
func ClientOptions(options ...*options.ClientOptions) Option {
	return func(cfg *Config) {
		cfg.ClientOptions = append(cfg.ClientOptions, options...)
	}
}

// EventStore ...
func EventStore(ctx context.Context, eventCfg cqrs.EventConfig, opts ...Option) (cqrs.EventStore, error) {
	var cfg Config
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.URI == "" {
		cfg.URI = os.Getenv("MONGO_EVENTS_URI")
	}

	if cfg.URI == "" {
		cfg.URI = os.Getenv("MONGO_URI")
	}

	if cfg.URI == "" {
		cfg.URI = DefaultURI
	}

	if cfg.Database == "" {
		cfg.Database = os.Getenv("MONGO_EVENTS_DB")
	}

	if cfg.Database == "" {
		cfg.Database = os.Getenv("MONGO_DB")
	}

	if cfg.Database == "" {
		cfg.Database = "events"
	}

	clientOptions := append([]*options.ClientOptions{options.Client().ApplyURI(cfg.URI)}, cfg.ClientOptions...)
	client, err := mongo.Connect(ctx, clientOptions...)
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.Database)

	if cfg.CreateIndexes {
		if err := createIndexes(ctx, db); err != nil {
			return nil, err
		}
	}

	return &eventStore{
		config:   cfg,
		db:       db,
		col:      db.Collection("events"),
		eventCfg: eventCfg,
	}, nil
}

// WithEventStoreFactory ...
func WithEventStoreFactory(options ...Option) setup.Option {
	return setup.WithEventStoreFactory(func(ctx context.Context, c cqrs.Container) (cqrs.EventStore, error) {
		options = append([]Option{ResolvePublisher(c.EventPublisher)}, options...)
		return EventStore(ctx, c.EventConfig(), options...)
	})
}

func (s *eventStore) Save(ctx context.Context, originalVersion int, events ...cqrs.Event) error {
	if len(events) == 0 {
		return nil
	}

	if err := event.Validate(events, originalVersion); err != nil {
		return fmt.Errorf("validate events: %w", err)
	}

	dbEvents := make([]*dbEvent, len(events))
	for i, e := range events {
		data := e.Data()
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(&data); err != nil {
			return fmt.Errorf("encode event data: %w", err)
		}

		dbevent := &dbEvent{
			EventType:     e.Type(),
			EventData:     buf.Bytes(),
			Time:          e.Time(),
			AggregateType: e.AggregateType(),
			AggregateID:   e.AggregateID(),
			Version:       e.Version(),
		}

		dbEvents[i] = dbevent
	}

	docs := make([]interface{}, len(dbEvents))
	for i, dbevent := range dbEvents {
		docs[i] = dbevent
	}

	if s.config.Transactions {
		if err := s.db.Client().UseSession(ctx, func(ctx mongo.SessionContext) error {
			if err := ctx.StartTransaction(); err != nil {
				return fmt.Errorf("start transaction: %w", err)
			}

			if err := s.saveDocs(ctx, events[0].AggregateType(), events[0].AggregateID(), originalVersion, docs); err != nil {
				return fmt.Errorf("save events: %w", err)
			}

			if err := ctx.CommitTransaction(ctx); err != nil {
				return fmt.Errorf("commit transaction: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}
	} else {
		if err := s.saveDocs(ctx, events[0].AggregateType(), events[0].AggregateID(), originalVersion, docs); err != nil {
			return fmt.Errorf("save events: %w", err)
		}
	}

	if err := s.publish(context.Background(), events...); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}

func (s *eventStore) publish(ctx context.Context, events ...cqrs.Event) error {
	if s.publisher == nil && s.config.ResolvePublisher != nil {
		if pub := s.config.ResolvePublisher(); pub != nil {
			s.publisher = pub
		}
	}

	if s.publisher != nil {
		return s.publisher.Publish(ctx, events...)
	}

	return nil
}

func (s *eventStore) saveDocs(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, originalVersion int, docs []interface{}) error {
	res := s.col.FindOne(ctx, bson.M{
		"aggregateType": aggregateType,
		"aggregateId":   aggregateID,
	}, options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}}))

	var latest dbEvent
	err := res.Decode(&latest)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("decode event: %w", err)
	}

	if err == nil && latest.Version != originalVersion {
		return event.OptimisticConcurrencyError{
			AggregateType:   aggregateType,
			AggregateID:     aggregateID,
			LatestVersion:   latest.Version,
			ProvidedVersion: originalVersion,
		}
	}

	if _, err := s.col.InsertMany(ctx, docs); err != nil {
		return fmt.Errorf("insert many: %w", err)
	}

	return nil
}

func (s *eventStore) Find(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, version int) (cqrs.Event, error) {
	res := s.col.FindOne(ctx, bson.M{
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
		return nil, err
	}

	return evt, nil
}

func (s *eventStore) Fetch(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, from int, to int) ([]cqrs.Event, error) {
	if from > to {
		return []cqrs.Event{}, nil
	}

	cur, err := s.col.Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
		{Key: "version", Value: bson.D{
			{Key: "$gte", Value: from},
			{Key: "$lte", Value: to},
		}},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, err
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, err
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, err
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) FetchAll(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID) ([]cqrs.Event, error) {
	cur, err := s.col.Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, err
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, err
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, err
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) FetchFrom(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, from int) ([]cqrs.Event, error) {
	cur, err := s.col.Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
		{Key: "version", Value: bson.D{
			{Key: "$gte", Value: from},
		}},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, err
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, err
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, err
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) FetchTo(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, to int) ([]cqrs.Event, error) {
	cur, err := s.col.Find(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
		{Key: "version", Value: bson.D{
			{Key: "$lte", Value: to},
		}},
	}, options.Find().SetSort(bson.D{{Key: "version", Value: 1}}))

	if err != nil {
		return nil, err
	}

	var dbevents []dbEvent
	if err := cur.All(ctx, &dbevents); err != nil {
		return nil, err
	}

	events := make([]cqrs.Event, len(dbevents))
	for i, dbevent := range dbevents {
		evt, err := s.toCQRSEvent(dbevent)
		if err != nil {
			return nil, err
		}
		events[i] = evt
	}

	return events, nil
}

func (s *eventStore) RemoveAggregate(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID) error {
	_, err := s.col.DeleteMany(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
		{Key: "aggregateId", Value: aggregateID},
	})
	return err
}

func (s *eventStore) RemoveAggregateType(ctx context.Context, aggregateType cqrs.AggregateType) error {
	_, err := s.col.DeleteMany(ctx, bson.D{
		{Key: "aggregateType", Value: aggregateType},
	})
	return err
}

func (s *eventStore) Query(ctx context.Context, query cqrs.EventQuery) (cqrs.EventCursor, error) {
	filter := bson.D{}

	if types := query.EventTypes(); len(types) > 0 {
		filter = append(filter, bson.E{Key: "type", Value: bson.D{
			{Key: "$in", Value: types},
		}})
	}

	if types := query.AggregateTypes(); len(types) > 0 {
		filter = append(filter, bson.E{Key: "aggregateType", Value: bson.D{
			{Key: "$in", Value: types},
		}})
	}

	if ids := query.AggregateIDs(); len(ids) > 0 {
		filter = append(filter, bson.E{Key: "aggregateId", Value: bson.D{
			{Key: "$in", Value: ids},
		}})
	}

	var versionOr []bson.D

	if versions := query.Versions(); len(versions) > 0 {
		versionOr = append(versionOr, bson.D{
			{Key: "version", Value: bson.D{{Key: "$in", Value: versions}}},
		})
	}

	if versionRanges := query.VersionRanges(); len(versionRanges) > 0 {
		for _, versionRange := range versionRanges {
			if versionRange[0] > versionRange[1] {
				versionRange[0], versionRange[1] = versionRange[1], versionRange[0]
			}

			versionOr = append(versionOr, bson.D{
				{Key: "version", Value: bson.D{
					{Key: "$gte", Value: versionRange[0]},
					{Key: "$lte", Value: versionRange[1]},
				}},
			})
		}
	}

	if len(versionOr) > 0 {
		filter = append(filter, bson.E{Key: "$or", Value: versionOr})
	}

	cur, err := s.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return newCursor(s, cur), nil
}

func (s *eventStore) toCQRSEvent(evt dbEvent) (cqrs.Event, error) {
	data, err := s.eventCfg.NewData(evt.EventType)
	if err != nil {
		return nil, err
	}

	if err := gob.NewDecoder(bytes.NewReader(evt.EventData)).Decode(&data); err != nil {
		return nil, err
	}

	return event.NewAggregateEventWithTime(evt.EventType, data, evt.Time, evt.AggregateType, evt.AggregateID, evt.Version), nil
}

type dbEvent struct {
	EventType     cqrs.EventType     `bson:"type"`
	EventData     []byte             `bson:"data"`
	Time          time.Time          `bson:"time"`
	AggregateType cqrs.AggregateType `bson:"aggregateType"`
	AggregateID   uuid.UUID          `bson:"aggregateId"`
	Version       int                `bson:"version"`
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("events").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "time", Value: 1}}},
		{Keys: bson.D{{Key: "aggregateType", Value: 1}}},
		{Keys: bson.D{{Key: "aggregateId", Value: 1}}},
		{Keys: bson.D{{Key: "version", Value: 1}}},
	})

	return err
}
