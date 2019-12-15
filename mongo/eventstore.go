package mongo

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/bounoable/cqrs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventStore struct {
	db        *mongo.Database
	eventCfg  cqrs.EventConfig
	publisher cqrs.EventPublisher
}

// NewEventStore ...
func NewEventStore(ctx context.Context, eventCfg cqrs.EventConfig, addr, dbname string, publisher cqrs.EventPublisher) (cqrs.EventStore, error) {
	if publisher == nil {
		return nil, wrapError(errors.New("event publisher cannot be nil"))
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		return nil, wrapError(err)
	}

	return &eventStore{
		db:        client.Database(dbname),
		eventCfg:  eventCfg,
		publisher: publisher,
	}, nil
}

// WithEventStoreFactory ...
func WithEventStoreFactory(ctx context.Context, addr, dbname string) cqrs.Option {
	return cqrs.WithEventStoreFactory(func(ctx context.Context, c cqrs.Core) (cqrs.EventStore, error) {
		return NewEventStore(ctx, c.EventConfig(), addr, dbname, c.EventBus())
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

	if err := s.publisher.Publish(context.Background(), events...); err != nil {
		return wrapError(err)
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
