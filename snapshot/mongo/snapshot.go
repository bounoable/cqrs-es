package mongo

import (
	"context"
	"errors"
	"os"

	"github.com/bounoable/cqrs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// DefaultURI ...
	DefaultURI = "mongodb://localhost:27017"
)

// Config ...
type Config struct {
	URI           string
	Database      string
	ClientOptions []*options.ClientOptions
}

// Option ...
type Option func(*Config)

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

// Database ...
func Database(name string) Option {
	return func(cfg *Config) {
		cfg.Database = name
	}
}

// URI ...
func URI(uri string) Option {
	return func(cfg *Config) {
		cfg.URI = uri
	}
}

// ClientOptions ...
func ClientOptions(options ...*options.ClientOptions) Option {
	return func(cfg *Config) {
		cfg.ClientOptions = append(cfg.ClientOptions, options...)
	}
}

// NewSnapshotRepository ...
func NewSnapshotRepository(ctx context.Context, aggregateConfig cqrs.AggregateConfig, opts ...Option) (cqrs.SnapshotRepository, error) {
	if aggregateConfig == nil {
		return nil, cqrs.SnapshotError{
			Err:       errors.New("aggregate config cannot be nil"),
			StoreName: "mongo",
		}
	}

	var cfg Config
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.URI == "" {
		cfg.URI = os.Getenv("MONGO_SNAPSHOTS_URI")
	}

	if cfg.URI == "" {
		cfg.URI = DefaultURI
	}

	if cfg.Database == "" {
		cfg.Database = os.Getenv("MONGO_SNAPSHOTS_DB")
	}

	if cfg.Database == "" {
		cfg.Database = "snapshots"
	}

	clientOptions := append([]*options.ClientOptions{options.Client().ApplyURI(cfg.URI)}, cfg.ClientOptions...)
	client, err := mongo.Connect(ctx, clientOptions...)
	if err != nil {
		return nil, cqrs.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return &snapshotRepository{
		db:              client.Database(cfg.Database),
		aggregateConfig: aggregateConfig,
	}, nil
}

// WithSnapshotRepositoryFactory ...
func WithSnapshotRepositoryFactory(options ...Option) cqrs.Option {
	return cqrs.WithSnapshotRepositoryFactory(func(ctx context.Context, c cqrs.Core) (cqrs.SnapshotRepository, error) {
		return NewSnapshotRepository(ctx, c.AggregateConfig(), options...)
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
