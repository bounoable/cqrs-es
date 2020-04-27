package event

import (
	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// MatchType ...
func MatchType(typ cqrs.EventType) cqrs.EventMatcher {
	return func(evt cqrs.Event) bool {
		return evt.Type() == typ
	}
}

// MatchAggregateType ...
func MatchAggregateType(typ cqrs.AggregateType) cqrs.EventMatcher {
	return func(evt cqrs.Event) bool {
		return evt.AggregateType() == typ
	}
}

// MatchAggregateID ...
func MatchAggregateID(id uuid.UUID) cqrs.EventMatcher {
	return func(evt cqrs.Event) bool {
		return evt.AggregateID() == id
	}
}

// MatchAggregate ...
func MatchAggregate(typ cqrs.AggregateType, id uuid.UUID) cqrs.EventMatcher {
	return func(evt cqrs.Event) bool {
		return MatchAggregateType(typ)(evt) && MatchAggregateID(id)(evt)
	}
}

// MatchVersion ...
func MatchVersion(version int) cqrs.EventMatcher {
	return func(evt cqrs.Event) bool {
		return evt.Version() == version
	}
}

// MatchMinVersion ...
func MatchMinVersion(minVersion int) cqrs.EventMatcher {
	return func(evt cqrs.Event) bool {
		return evt.Version() >= minVersion
	}
}

// MatchMaxVersion ...
func MatchMaxVersion(maxVersion int) cqrs.EventMatcher {
	return func(evt cqrs.Event) bool {
		return evt.Version() <= maxVersion
	}
}
