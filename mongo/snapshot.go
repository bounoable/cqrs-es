package mongo

import (
	"context"
	"errors"

	"github.com/bounoable/cqrs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type snapshotRepository struct {
	db              *mongo.Database
	aggregateConfig cqrs.AggregateConfig
}

type snapshot struct {
	AggregateType cqrs.AggregateType `bson:"aggregateType"`
	AggregateID   uuid.UUID          `bson:"aggregateId"`
	Version       int                `bson:"version"`
	Data          []byte             `bson:"data"`
}

// NewSnapshotRepository ...
func NewSnapshotRepository(ctx context.Context, aggregateConfig cqrs.AggregateConfig, addr, dbname string) (cqrs.SnapshotRepository, error) {
	if aggregateConfig == nil {
		return nil, cqrs.SnapshotError{
			Err:       errors.New("aggregate config cannot be nil"),
			StoreName: "mongo",
		}
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return &snapshotRepository{
		db:              client.Database(dbname),
		aggregateConfig: aggregateConfig,
	}, nil
}

// WithSnapshotRepositoryFactory ...
func WithSnapshotRepositoryFactory(addr, dbname string) cqrs.Option {
	return cqrs.WithSnapshotRepositoryFactory(func(ctx context.Context, c cqrs.Core) (cqrs.SnapshotRepository, error) {
		return NewSnapshotRepository(ctx, c.AggregateConfig(), addr, dbname)
	})
}

func (r *snapshotRepository) Save(ctx context.Context, aggregate cqrs.Aggregate) error {
	data, err := bson.Marshal(aggregate)
	if err != nil {
		return cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	snap := &snapshot{
		AggregateType: aggregate.AggregateType(),
		AggregateID:   aggregate.AggregateID(),
		Version:       aggregate.CurrentVersion(),
		Data:          data,
	}

	if _, err := r.db.Collection("snapshots").ReplaceOne(ctx, bson.M{
		"aggregateType": snap.AggregateType,
		"aggregateId":   snap.AggregateID,
		"version":       snap.Version,
	}, snap, options.Replace().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (r *snapshotRepository) Find(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, version int) (cqrs.Aggregate, error) {
	res := r.db.Collection("snapshots").FindOne(ctx, bson.M{
		"aggregateType": typ,
		"aggregateId":   id,
		"version":       version,
	})

	var snap snapshot
	if err := res.Decode(&snap); err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	aggregate, err := r.aggregateConfig.New(snap.AggregateType, snap.AggregateID)
	if err != nil {
		return nil, err
	}

	if err := bson.Unmarshal(snap.Data, aggregate); err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return aggregate, nil
}

func (r *snapshotRepository) Latest(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error) {
	res := r.db.Collection("snapshots").FindOne(ctx, bson.M{
		"aggregateType": typ,
		"aggregateId":   id,
	}, options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}}))

	var snap snapshot
	if err := res.Decode(&snap); err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	aggregate, err := r.aggregateConfig.New(snap.AggregateType, snap.AggregateID)
	if err != nil {
		return nil, err
	}

	if err := bson.Unmarshal(snap.Data, aggregate); err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return aggregate, nil
}

func (r *snapshotRepository) MaxVersion(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, maxVersion int) (cqrs.Aggregate, error) {
	res := r.db.Collection("snapshots").FindOne(ctx, bson.D{
		{Key: "aggregateType", Value: typ},
		{Key: "aggregateId", Value: id},
		{Key: "version", Value: bson.D{
			{Key: "$lte", Value: maxVersion},
		}},
	}, options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}}))

	var snap snapshot
	if err := res.Decode(&snap); err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	aggregate, err := r.aggregateConfig.New(snap.AggregateType, snap.AggregateID)
	if err != nil {
		return nil, err
	}

	if err := bson.Unmarshal(snap.Data, aggregate); err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return aggregate, nil
}
