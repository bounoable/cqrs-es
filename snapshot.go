package cqrs

//go:generate mockgen -source=snapshot.go -destination=./mocks/snapshot.go

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SnapshotRepository ...
type SnapshotRepository interface {
	Save(ctx context.Context, snap Aggregate) error
	Find(ctx context.Context, typ AggregateType, id uuid.UUID, version int) (Aggregate, error)
	Latest(ctx context.Context, typ AggregateType, id uuid.UUID) (Aggregate, error)
	MaxVersion(ctx context.Context, typ AggregateType, id uuid.UUID, maxVersion int) (Aggregate, error)
}

// SnapshotConfig ...
type SnapshotConfig interface {
	IsDue(Aggregate) bool
}

type snapshotConfig struct {
	every int
}

// SnapshotError ...
type SnapshotError struct {
	Err       error
	StoreName string
}

func (err SnapshotError) Error() string {
	return fmt.Sprintf("%s snapshot: %s", err.StoreName, err.Err)
}

// Unwrap ...
func (err SnapshotError) Unwrap() error {
	return err.Err
}

// NewSnapshotConfig ...
func NewSnapshotConfig(every int) SnapshotConfig {
	return snapshotConfig{
		every: every,
	}
}

func (cfg snapshotConfig) IsDue(aggregate Aggregate) bool {
	return cfg.every > 0 && len(aggregate.Changes()) >= cfg.every
}
