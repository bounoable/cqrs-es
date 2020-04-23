package aggregate

import (
	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// QueryOption ...
type QueryOption func(*query)

// QueryType ...
func QueryType(types ...cqrs.AggregateType) QueryOption {
	return func(q *query) {
		q.types = types
	}
}

// QueryID ...
func QueryID(ids ...uuid.UUID) QueryOption {
	return func(q *query) {
		q.ids = ids
	}
}

// Query ...
func Query(opts ...QueryOption) cqrs.AggregateQuery {
	var q query
	for _, opt := range opts {
		opt(&q)
	}
	return q
}

type query struct {
	types []cqrs.AggregateType
	ids   []uuid.UUID
}

func (q query) Types() []cqrs.AggregateType {
	return q.types
}

func (q query) IDs() []uuid.UUID {
	return q.ids
}
