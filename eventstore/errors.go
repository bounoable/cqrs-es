package eventstore

import (
	"fmt"

	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// EventNotFoundError ...
type EventNotFoundError struct {
	AggregateType cqrs.AggregateType
	AggregateID   uuid.UUID
	Version       int
}

func (err EventNotFoundError) Error() string {
	return fmt.Sprintf("event not found for %s:%s@%d", err.AggregateType, err.AggregateID, err.Version)
}
