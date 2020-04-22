package event

import (
	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// QueryOption ...
type QueryOption func(*query)

// QueryEventType ...
func QueryEventType(types ...cqrs.EventType) QueryOption {
	return func(q *query) {
		if types == nil {
			types = []cqrs.EventType{}
		}
		q.eventTypes = types
	}
}

// QueryAggregateType ...
func QueryAggregateType(types ...cqrs.AggregateType) QueryOption {
	return func(q *query) {
		if types == nil {
			types = []cqrs.AggregateType{}
		}
		q.aggregateTypes = types
	}
}

// QueryAggregateID ...
func QueryAggregateID(ids ...uuid.UUID) QueryOption {
	return func(q *query) {
		if ids == nil {
			ids = []uuid.UUID{}
		}
		q.aggregateIDs = ids
	}
}

// QueryVersions ...
func QueryVersions(versions ...int) QueryOption {
	return func(q *query) {
		if versions == nil {
			versions = []int{}
		}
		q.versions = versions
	}
}

// QueryVersionRanges ...
func QueryVersionRanges(ranges ...[2]int) QueryOption {
	return func(q *query) {
		if ranges == nil {
			ranges = [][2]int{}
		}
		q.versionRanges = ranges
	}
}

// Query ...
func Query(opts ...QueryOption) cqrs.EventQuery {
	var q query
	for _, opt := range opts {
		opt(&q)
	}
	return q
}

type query struct {
	eventTypes     []cqrs.EventType
	aggregateTypes []cqrs.AggregateType
	aggregateIDs   []uuid.UUID
	versions       []int
	versionRanges  [][2]int
}

func (q query) EventTypes() []cqrs.EventType {
	return q.eventTypes
}

func (q query) AggregateTypes() []cqrs.AggregateType {
	return q.aggregateTypes
}

func (q query) AggregateIDs() []uuid.UUID {
	return q.aggregateIDs
}

func (q query) Versions() []int {
	return q.versions
}

func (q query) VersionRanges() [][2]int {
	return q.versionRanges
}
