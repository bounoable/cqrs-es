package aggregate

import (
	"context"

	"github.com/bounoable/cqrs-es"
)

type cursor struct {
	aggregates   []cqrs.Aggregate
	currentIndex int
	err          error
}

func newCursor(aggregates []cqrs.Aggregate) cursor {
	return cursor{
		aggregates:   aggregates,
		currentIndex: -1,
	}
}

func (cur cursor) Next(_ context.Context) bool {
	if cur.currentIndex >= len(cur.aggregates)-1 {
		return false
	}
	cur.currentIndex++
	return true
}

func (cur cursor) Aggregate() cqrs.Aggregate {
	return cur.aggregates[cur.currentIndex]
}

func (cur cursor) Err() error {
	return cur.err
}

func (cur cursor) Close(_ context.Context) error {
	return nil
}
