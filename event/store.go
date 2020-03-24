package event

import (
	"fmt"

	cqrs "github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// StoreError ...
type StoreError struct {
	Err       error
	StoreName string
}

func (err StoreError) Error() string {
	return fmt.Sprintf("%s event store: %s", err.StoreName, err.Err)
}

// Unwrap ...
func (err StoreError) Unwrap() error {
	return err.Err
}

// OptimisticConcurrencyError ...
type OptimisticConcurrencyError struct {
	AggregateType   cqrs.AggregateType
	AggregateID     uuid.UUID
	LatestVersion   int
	ProvidedVersion int
}

func (err OptimisticConcurrencyError) Error() string {
	return fmt.Sprintf(
		"optimistic concurrency (%s:%s): latest version is %v, provided version is %v",
		err.AggregateType,
		err.AggregateID,
		err.LatestVersion,
		err.ProvidedVersion,
	)
}
