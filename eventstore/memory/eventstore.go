package memory

import (
	"context"
	"sync"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/event"
	"github.com/google/uuid"
)

type eventStore struct {
	mux      sync.RWMutex
	events   map[cqrs.AggregateType]map[uuid.UUID][]cqrs.Event
	eventPub cqrs.EventPublisher
}

// Option ...
type Option func(*eventStore)

// EventPublisher ...
func EventPublisher(eventPub cqrs.EventPublisher) Option {
	return func(s *eventStore) {
		s.eventPub = eventPub
	}
}

// EventStore ...
func EventStore(opts ...Option) cqrs.EventStore {
	s := &eventStore{
		events: make(map[cqrs.AggregateType]map[uuid.UUID][]cqrs.Event),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *eventStore) Save(ctx context.Context, originalVersion int, events ...cqrs.Event) error {
	if len(events) == 0 {
		return nil
	}

	if err := event.Validate(events, originalVersion); err != nil {
		return err
	}

	aggregateType := events[0].AggregateType()
	aggregateID := events[0].AggregateID()

	s.mux.Lock()
	if _, ok := s.events[aggregateType]; !ok {
		s.events[aggregateType] = make(map[uuid.UUID][]cqrs.Event)
	}

	if _, ok := s.events[aggregateType][aggregateID]; !ok {
		s.events[aggregateType][aggregateID] = make([]cqrs.Event, 0)
	}
	s.events[aggregateType][aggregateID] = append(s.events[aggregateType][aggregateID], events...)
	s.mux.Unlock()

	return s.eventPub.Publish(ctx, events...)
}

func (s *eventStore) Find(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, version int) (cqrs.Event, error) {
	events := s.aggregateEvents(aggregateType, aggregateID)

	for _, evt := range events {
		if evt.Version() == version {
			return evt, nil
		}
	}

	return nil, event.NotFoundError{
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Version:       version,
	}
}

func (s *eventStore) Fetch(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, from int, to int) ([]cqrs.Event, error) {
	events := s.aggregateEvents(aggregateType, aggregateID)

	filtered := []cqrs.Event{}
	for _, evt := range events {
		if evt.Version() < from || evt.Version() > to {
			continue
		}
		filtered = append(filtered, evt)
	}

	return filtered, nil
}

func (s *eventStore) FetchAll(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID) ([]cqrs.Event, error) {
	return s.aggregateEvents(aggregateType, aggregateID), nil
}

func (s *eventStore) FetchFrom(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, from int) ([]cqrs.Event, error) {
	events := s.aggregateEvents(aggregateType, aggregateID)

	filtered := []cqrs.Event{}
	for _, evt := range events {
		if evt.Version() < from {
			continue
		}
		filtered = append(filtered, evt)
	}

	return filtered, nil
}

func (s *eventStore) FetchTo(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID, to int) ([]cqrs.Event, error) {
	events := s.aggregateEvents(aggregateType, aggregateID)

	filtered := []cqrs.Event{}
	for _, evt := range events {
		if evt.Version() > to {
			continue
		}
		filtered = append(filtered, evt)
	}

	return filtered, nil
}

func (s *eventStore) RemoveAll(ctx context.Context, aggregateType cqrs.AggregateType, aggregateID uuid.UUID) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.events[aggregateType]; ok {
		delete(s.events[aggregateType], aggregateID)
	}

	return nil
}

func (s *eventStore) aggregateEvents(typ cqrs.AggregateType, id uuid.UUID) []cqrs.Event {
	s.mux.RLock()
	defer s.mux.RUnlock()
	typeEvents, ok := s.events[typ]
	if !ok {
		return []cqrs.Event{}
	}

	events, ok := typeEvents[id]
	if !ok {
		return []cqrs.Event{}
	}

	return events
}
