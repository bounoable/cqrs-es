package aggregate

import (
	"context"

	"github.com/bounoable/cqrs-es"
	"github.com/google/uuid"
)

// Use ...
func Use(
	ctx context.Context,
	repo cqrs.AggregateRepository,
	typ cqrs.AggregateType,
	id uuid.UUID,
	version int,
	use func(context.Context, cqrs.Aggregate) error,
) (cqrs.Aggregate, error) {
	aggregate, err := repo.Fetch(ctx, typ, id, version)
	if err != nil {
		return nil, err
	}

	return aggregate, use(ctx, aggregate)
}

// UseLatest ...
func UseLatest(
	ctx context.Context,
	repo cqrs.AggregateRepository,
	typ cqrs.AggregateType,
	id uuid.UUID,
	use func(context.Context, cqrs.Aggregate) error,
) (cqrs.Aggregate, error) {
	aggregate, err := repo.FetchLatest(ctx, typ, id)
	if err != nil {
		return nil, err
	}

	return aggregate, use(ctx, aggregate)
}
