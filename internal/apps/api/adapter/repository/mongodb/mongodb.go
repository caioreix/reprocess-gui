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

func (db *db) Close(ctx context.Context) error {
	return db.Disconnect(ctx)
}
