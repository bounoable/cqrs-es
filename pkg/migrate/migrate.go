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

	buf := make(chan cqrs.Event, 100)

	var iter int
	for cur.Next(ctx) {
		select {
		case buf <- cur.Event():
		default:
			close(buf)
			if err := m.insertBuffer(ctx, buf, iter); err != nil {
				return fmt.Errorf("insert buffer: %w", err)
			}
			buf = make(chan cqrs.Event, 100)
			buf <- cur.Event()
			iter++
		}
	}
	close(buf)

	if err := m.insertBuffer(ctx, buf, iter); err != nil {
		return fmt.Errorf("insert buffer: %w", err)
	}

	if err := cur.Err(); err != nil {
		return fmt.Errorf("cursor: %w", err)
	}

	return nil
}

func (m *Migrator) insertBuffer(ctx context.Context, buf <-chan cqrs.Event, iter int) error {
	if len(buf) == 0 {
		return nil
	}
	start := 100 * iter
	end := start + len(buf) - 1
	log.Printf("Inserting Events %d-%d...\n", start, end)
	var events []event.Event
	for evt := range buf {
		events = append(events, event.New(
			evt.Type().String(),
			evt.Data(),
			event.Time(evt.Time()),
			event.Aggregate(string(evt.AggregateType()), evt.AggregateID(), evt.Version()+1),
		))
	}
	return m.store.Insert(ctx, events...)
}
