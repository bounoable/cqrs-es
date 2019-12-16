package cqrs

//go:generate mockgen -source=snapshot.go -destination=./mocks/snapshot.go

import (
	"context"
	"fmt"
	"sync"

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

// SnapshotOption ...
type SnapshotOption func(*snapshotConfig)

type snapshotConfig struct {
	mux       sync.RWMutex
	intervals map[AggregateType]int
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

// SnapshotInterval ...
func SnapshotInterval(typ AggregateType, every int) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.intervals[typ] = every
	}
}

// NewSnapshotConfig ...
func NewSnapshotConfig(options ...SnapshotOption) SnapshotConfig {
	var cfg snapshotConfig
	for _, opt := range options {
		opt(&cfg)
	}

	return &cfg
}

func (cfg *snapshotConfig) IsDue(aggregate Aggregate) bool {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	interval, ok := cfg.intervals[aggregate.AggregateType()]
	if !ok {
		return false
	}

	return interval > 0 && len(aggregate.Changes()) >= interval
}
