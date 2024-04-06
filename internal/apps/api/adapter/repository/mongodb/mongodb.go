package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"reprocess-gui/internal/apps/api/config"
)

type db struct {
	*mongo.Client
}

// New creates a new MongoDB client based on the provided configuration.
// It establishes a connection to the MongoDB server and performs a ping to ensure connectivity.
func New(config *config.Config) (*db, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Mongo.ConnectionTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.URI))
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.Mongo.PingTimeout)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &db{
		Client: client,
	}, nil
}

// Close disconnects the MongoDB client gracefully.
func (db *db) Close(ctx context.Context) error {
	return db.Disconnect(ctx)
}
