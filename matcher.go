package cqrs

import "github.com/google/uuid"

// Matcher ...
type Matcher func(Event) bool

// MatchEventType ...
func MatchEventType(typ EventType) Matcher {
	return func(evt Event) bool {
		return evt.Type() == typ
	}
}

// MatchAggregateType ...
func MatchAggregateType(typ AggregateType) Matcher {
	return func(evt Event) bool {
		return evt.AggregateType() == typ
	}
}

// MatchAggregateID ...
func MatchAggregateID(id uuid.UUID) Matcher {
	return func(evt Event) bool {
		return evt.AggregateID() == id
	}
}

// MatchAggregate ...
func MatchAggregate(typ AggregateType, id uuid.UUID) Matcher {
	return func(evt Event) bool {
		return MatchAggregateType(typ)(evt) && MatchAggregateID(id)(evt)
	}
}

// MatchVersion ...
func MatchVersion(version int) Matcher {
	return func(evt Event) bool {
		return evt.Version() == version
	}
}

// MatchMinVersion ...
func MatchMinVersion(minVersion int) Matcher {
	return func(evt Event) bool {
		return evt.Version() >= minVersion
	}
}

// MatchMaxVersion ...
func MatchMaxVersion(maxVersion int) Matcher {
	return func(evt Event) bool {
		return evt.Version() <= maxVersion
	}
}
