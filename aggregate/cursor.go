package aggregate

import (
	"context"
	"errors"

	"github.com/bounoable/cqrs-es"
)

type cursor struct {
	aggregates cqrs.AggregateRepository
	// eventcursor should only return 0-version events
	eventcursor cqrs.EventCursor
	err         error
}

func newCursor(aggregates cqrs.AggregateRepository, eventcursor cqrs.EventCursor) cursor {
	return cursor{
		aggregates:  aggregates,
		eventcursor: eventcursor,
	}
}

func (cur cursor) Next(ctx context.Context) bool {
	if !cur.eventcursor.Next(ctx) {
		cur.err = cur.eventcursor.Err()
		return false
	}
	return true
}

func (cur cursor) Aggregate(ctx context.Context) (cqrs.Aggregate, error) {
	evt, err := cur.currentEvent()
	if err != nil {
		return nil, err
	}

	agg, err := cur.aggregates.FetchLatest(ctx, evt.AggregateType(), evt.AggregateID())
	if err != nil {
		return nil, err
	}

	return agg, nil
}

func (cur cursor) Version(ctx context.Context, version int) (cqrs.Aggregate, error) {
	evt, err := cur.currentEvent()
	if err != nil {
		return nil, err
	}

	agg, err := cur.aggregates.Fetch(ctx, evt.AggregateType(), evt.AggregateID(), version)
	if err != nil {
		return nil, err
	}

	return agg, nil
}

func (cur cursor) currentEvent() (cqrs.Event, error) {
	evt := cur.eventcursor.Event()
	if evt == nil {
		return nil, errors.New("cursor has no current event")
	}

	return evt, nil
}

func (cur cursor) Err() error {
	return cur.err
}

func (cur cursor) Close(ctx context.Context) error {
	return cur.eventcursor.Close(ctx)
}
