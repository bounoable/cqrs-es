package aggregate

import (
	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// MatchType ...
func MatchType(types ...cqrs.AggregateType) cqrs.AggregateMatcher {
	return func(agg cqrs.Aggregate) bool {
		for _, typ := range types {
			if typ == agg.AggregateType() {
				return true
			}
		}
		return false
	}
}

// MatchID ...
func MatchID(ids ...uuid.UUID) cqrs.AggregateMatcher {
	return func(agg cqrs.Aggregate) bool {
		for _, id := range ids {
			if id == agg.AggregateID() {
				return true
			}
		}
		return false
	}
}

// MatchVersion ...
func MatchVersion(versions ...int) cqrs.AggregateMatcher {
	return func(agg cqrs.Aggregate) bool {
		for _, version := range versions {
			if version == agg.CurrentVersion() {
				return true
			}
		}
		return false
	}
}

// MatchMinVersion ...
func MatchMinVersion(version int) cqrs.AggregateMatcher {
	return func(agg cqrs.Aggregate) bool {
		if agg.CurrentVersion() >= version {
			return true
		}
		return false
	}
}

// MatchMaxVersion ...
func MatchMaxVersion(version int) cqrs.AggregateMatcher {
	return func(agg cqrs.Aggregate) bool {
		if agg.CurrentVersion() <= version {
			return true
		}
		return false
	}
}
