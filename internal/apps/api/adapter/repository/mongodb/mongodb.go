package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func New(ctx context.Context, config *config.Config) (*db, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Mongo.ConnectionTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.URI))
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(ctx, config.Mongo.PingTimeout)
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

// GetTotalCount returns the total count of documents in the specified collection using the database collStats.
func GetTotalCount(ctx context.Context, coll *mongo.Collection) (int, error) {
	sr := coll.Database().RunCommand(ctx, bson.M{"collStats": coll.Name()})
	result := struct {
		Count int
	}{}
	err := sr.Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

// buildPaginationFilterAndOptions constructs the filter and options for the MongoDB query.
func buildPaginationFilterAndOptions(key, offset string, limit int, reversed bool) (any, *options.FindOptions, error) {
	var filter any = bson.D{}
	opts := options.Find().SetLimit(int64(limit))

	if offset != "" {
		id, err := primitive.ObjectIDFromHex(offset)
		if err != nil {
			return nil, nil, err
		}
		filter = bson.M{key: bson.M{"$gt": id}}
	}

	if reversed {
		opts = opts.SetSort(bson.M{key: -1})
		if offset != "" {
			id, _ := primitive.ObjectIDFromHex(offset)
			filter = bson.M{key: bson.M{"$lt": id}}
		}
	}

	return filter, opts, nil
}

// decodePaginationResponse decodes the MongoDB cursor into a slice of T reversing the cursor if activated.
func decodePaginationResponse[T any](ctx context.Context, cur *mongo.Cursor, reversed bool) ([]*T, error) {
	items := []*T{}

	for cur.Next(ctx) {
		var item T
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}

		if reversed {
			items = append([]*T{&item}, items...)
		} else {
			items = append(items, &item)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
