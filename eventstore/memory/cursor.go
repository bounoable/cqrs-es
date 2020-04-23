package memory

import (
	"context"

	"github.com/bounoable/cqrs-es"
)

type cursor struct {
	events       []cqrs.Event
	currentIndex int
}

func newCursor(events []cqrs.Event) *cursor {
	return &cursor{
		events:       events,
		currentIndex: -1,
	}
}

func (cur *cursor) Next(ctx context.Context) bool {
	if cur.currentIndex >= len(cur.events)-1 {
		return false
	}
	cur.currentIndex++
	return true
}

func (cur *cursor) Event() cqrs.Event {
	return cur.events[cur.currentIndex]
}

func (cur *cursor) Err() error {
	return nil
}

func (cur *cursor) Close(_ context.Context) error {
	return nil
}
