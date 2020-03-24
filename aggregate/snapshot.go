package aggregate

//go:generate mockgen -source=snapshot.go -destination=../mocks/aggregate/snapshot.go

import (
	"fmt"
	"sync"

	cqrs "github.com/bounoable/cqrs-es"
)

// SnapshotConfig ...
type SnapshotConfig interface {
	IsDue(cqrs.Aggregate) bool
}

type snapshotConfig struct {
	mux       sync.RWMutex
	intervals map[cqrs.AggregateType]int
}

// SnapshotOption ...
type SnapshotOption func(*snapshotConfig)

// SnapshotInterval ...
func SnapshotInterval(typ cqrs.AggregateType, every int) SnapshotOption {
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

func (cfg *snapshotConfig) IsDue(aggregate cqrs.Aggregate) bool {
	cfg.mux.RLock()
	defer cfg.mux.RUnlock()

	interval, ok := cfg.intervals[aggregate.AggregateType()]
	if !ok {
		return false
	}

	return interval > 0 && len(aggregate.Changes()) >= interval
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
