package mongo

import (
	"context"
	"errors"
	"os"

	"github.com/bounoable/cqrs-es"
	"github.com/bounoable/cqrs-es/aggregate"
	"github.com/bounoable/cqrs-es/setup"
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
	CreateIndexes bool
	ClientOptions []*options.ClientOptions
}

// Option ...
type Option func(*Config)

type snapshotRepository struct {
	config          Config
	col             *mongo.Collection
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

// CreateIndexes ...
func CreateIndexes() Option {
	return func(cfg *Config) {
		cfg.CreateIndexes = true
	}
}

// ClientOptions ...
func ClientOptions(options ...*options.ClientOptions) Option {
	return func(cfg *Config) {
		cfg.ClientOptions = append(cfg.ClientOptions, options...)
	}
}

// SnapshotRepository ...
func SnapshotRepository(ctx context.Context, aggregateConfig cqrs.AggregateConfig, opts ...Option) (cqrs.SnapshotRepository, error) {
	if aggregateConfig == nil {
		return nil, aggregate.SnapshotError{
			Err:       errors.New("nil aggregate config"),
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
		return nil, aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	db := client.Database(cfg.Database)

	if cfg.CreateIndexes {
		if err := createIndexes(ctx, db); err != nil {
			return nil, aggregate.SnapshotError{
				Err:       err,
				StoreName: "mongo",
			}
		}
	}

	return &snapshotRepository{
		config:          cfg,
		col:             db.Collection("snapshots"),
		aggregateConfig: aggregateConfig,
	}, nil
}

// WithSnapshotRepositoryFactory ...
func WithSnapshotRepositoryFactory(options ...Option) setup.Option {
	return setup.WithSnapshotRepositoryFactory(func(ctx context.Context, c cqrs.Container) (cqrs.SnapshotRepository, error) {
		return SnapshotRepository(ctx, c.AggregateConfig(), options...)
	})
}

func (r *snapshotRepository) Save(ctx context.Context, agg cqrs.Aggregate) error {
	data, err := bson.Marshal(agg)
	if err != nil {
		return aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	snap := &snapshot{
		AggregateType: agg.AggregateType(),
		AggregateID:   agg.AggregateID(),
		Version:       agg.CurrentVersion(),
		Data:          data,
	}

	if _, err := r.col.ReplaceOne(ctx, bson.M{
		"aggregateType": snap.AggregateType,
		"aggregateId":   snap.AggregateID,
		"version":       snap.Version,
	}, snap, options.Replace().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (r *snapshotRepository) Find(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, version int) (cqrs.Aggregate, error) {
	res := r.col.FindOne(ctx, bson.M{
		"aggregateType": typ,
		"aggregateId":   id,
		"version":       version,
	})

	var snap snapshot
	if err := res.Decode(&snap); err != nil {
		return nil, aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	agg, err := r.aggregateConfig.New(snap.AggregateType, snap.AggregateID)
	if err != nil {
		return nil, err
	}

	if err := bson.Unmarshal(snap.Data, agg); err != nil {
		return nil, aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return agg, nil
}

func (r *snapshotRepository) Latest(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) (cqrs.Aggregate, error) {
	res := r.col.FindOne(ctx, bson.M{
		"aggregateType": typ,
		"aggregateId":   id,
	}, options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}}))

	var snap snapshot
	if err := res.Decode(&snap); err != nil {
		return nil, aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	agg, err := r.aggregateConfig.New(snap.AggregateType, snap.AggregateID)
	if err != nil {
		return nil, err
	}

	if err := bson.Unmarshal(snap.Data, agg); err != nil {
		return nil, aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return agg, nil
}

func (r *snapshotRepository) MaxVersion(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, maxVersion int) (cqrs.Aggregate, error) {
	res := r.col.FindOne(ctx, bson.D{
		{Key: "aggregateType", Value: typ},
		{Key: "aggregateId", Value: id},
		{Key: "version", Value: bson.D{
			{Key: "$lte", Value: maxVersion},
		}},
	}, options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}}))

	var snap snapshot
	if err := res.Decode(&snap); err != nil {
		return nil, aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	agg, err := r.aggregateConfig.New(snap.AggregateType, snap.AggregateID)
	if err != nil {
		return nil, err
	}

	if err := bson.Unmarshal(snap.Data, agg); err != nil {
		return nil, aggregate.SnapshotError{
			Err:       err,
			StoreName: "mongo",
		}
	}

	return agg, nil
}

func (r *snapshotRepository) Remove(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID, version int) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"aggregateType": typ, "aggregateId": id, "version": version})
	return err
}

func (r *snapshotRepository) RemoveAll(ctx context.Context, typ cqrs.AggregateType, id uuid.UUID) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"aggregateType": typ, "aggregateId": id})
	return err
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("snapshots").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "aggregateType", Value: 1}}},
		{Keys: bson.D{{Key: "aggregateId", Value: 1}}},
		{Keys: bson.D{{Key: "version", Value: 1}}},
	})

	return err
}
