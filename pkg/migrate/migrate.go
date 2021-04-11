package migrate

import (
	"context"
	"fmt"
	"log"

	"github.com/bounoable/cqrs-es"
	cqrsevent "github.com/bounoable/cqrs-es/event"
	"github.com/modernice/goes/event"
)

type Migrator struct {
	oldstore cqrs.EventStore
	store    event.Store
}

func NewMigrator(
	oldstore cqrs.EventStore,
	store event.Store,
) *Migrator {
	return &Migrator{
		oldstore: oldstore,
		store:    store,
	}
}

func (m *Migrator) Migrate(ctx context.Context) error {
	cur, err := m.oldstore.Query(ctx, cqrsevent.Query())
	if err != nil {
		return fmt.Errorf("query cqrs Events: %w", err)
	}
	defer cur.Close(ctx)

	buf := make([]cqrs.Event, 0, 1000)

	var iter int
	for cur.Next(ctx) {
		buf = append(buf, cur.Event())
		if len(buf) >= 1000 {
			if err := m.insertBuffer(ctx, buf, iter); err != nil {
				return fmt.Errorf("insert buffer: %w", err)
			}
			buf = buf[:0]
			iter++
		}
	}

	if err := m.insertBuffer(ctx, buf, iter); err != nil {
		return fmt.Errorf("insert buffer: %w", err)
	}

	if err := cur.Err(); err != nil {
		return fmt.Errorf("cursor: %w", err)
	}

	return nil
}

func (m *Migrator) insertBuffer(ctx context.Context, buf []cqrs.Event, iter int) error {
	if len(buf) == 0 {
		return nil
	}
	start := 1000 * iter
	end := start + len(buf) - 1
	log.Printf("Inserting Events %d-%d...\n", start, end)
	var events []event.Event
	for _, evt := range buf {
		events = append(events, event.New(
			evt.Type().String(),
			evt.Data(),
			event.Time(evt.Time()),
			event.Aggregate(string(evt.AggregateType()), evt.AggregateID(), evt.Version()+1),
		))
	}
	return m.store.Insert(ctx, events...)
}
