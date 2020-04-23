package mongo

import (
	"context"

	"github.com/bounoable/cqrs-es"
	"go.mongodb.org/mongo-driver/mongo"
)

type cursor struct {
	store   *eventStore
	cursor  *mongo.Cursor
	current cqrs.Event
	err     error
}

func newCursor(store *eventStore, cur *mongo.Cursor) *cursor {
	return &cursor{
		store:  store,
		cursor: cur,
	}
}

func (cur *cursor) Next(ctx context.Context) bool {
	if !cur.cursor.Next(ctx) {
		cur.err = cur.cursor.Err()
		return false
	}

	var dbevent dbEvent
	err := cur.cursor.Decode(&dbevent)
	if err != nil {
		cur.err = err
		return false
	}

	if cur.current, cur.err = cur.store.toCQRSEvent(dbevent); cur.err != nil {
		return false
	}

	return true
}

func (cur *cursor) Event() cqrs.Event {
	return cur.current
}

func (cur *cursor) Err() error {
	return cur.err
}

func (cur *cursor) Close(ctx context.Context) error {
	return cur.cursor.Close(ctx)
}
