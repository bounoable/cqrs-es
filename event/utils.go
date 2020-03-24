package event

import (
	"fmt"

	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// Validate validates events.
func Validate(events []cqrs.Event, originalVersion int) error {
	if len(events) == 0 {
		return nil
	}

	aggregateType := events[0].AggregateType()
	aggregateID := events[0].AggregateID()

	for i, evt := range events {
		if evt.AggregateType() != aggregateType || evt.AggregateID() != aggregateID {
			return AggregateMismatchError{
				ExpectedAggregateType: aggregateType,
				ProvidedAggregateType: evt.AggregateType(),
				ExpectedAggregateID:   aggregateID,
				ProvidedAggregateID:   evt.AggregateID(),
			}
		}

		if evt.Version() != originalVersion+i+1 {
			return InconsistentVersionError{
				OriginalVersion: originalVersion,
				ExpectedVersion: originalVersion + i + 1,
				ProvidedVersion: evt.Version(),
			}
		}
	}

	return nil
}

// AggregateMismatchError ...
type AggregateMismatchError struct {
	ExpectedAggregateType cqrs.AggregateType
	ProvidedAggregateType cqrs.AggregateType
	ExpectedAggregateID   uuid.UUID
	ProvidedAggregateID   uuid.UUID
}

func (err AggregateMismatchError) Error() string {
	return fmt.Sprintf(
		"aggregate mismatch: expected '%s:%s', got '%s:%s'",
		err.ExpectedAggregateType,
		err.ExpectedAggregateID,
		err.ProvidedAggregateType,
		err.ProvidedAggregateID,
	)
}

// InconsistentVersionError ...
type InconsistentVersionError struct {
	OriginalVersion int
	ExpectedVersion int
	ProvidedVersion int
}

func (err InconsistentVersionError) Error() string {
	return fmt.Sprintf("inconsistent version: expected %d, got %d", err.ExpectedVersion, err.ProvidedVersion)
}
