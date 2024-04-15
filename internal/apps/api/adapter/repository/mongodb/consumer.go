package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/logger"
	"reprocess-gui/internal/utils"
)

type consumerRepository struct {
	config     *config.Config
	log        *logger.Logger
	collection *mongo.Collection
}

// NewConsumerRepository creates a new instance of consumerRepository.
func NewConsumerRepository(config *config.Config, log *logger.Logger, collection *mongo.Collection) *consumerRepository {
	return &consumerRepository{
		config:     config,
		log:        log,
		collection: collection,
	}
}

// GetAllConsumers retrieves all consumers from the MongoDB collection.
func (r *consumerRepository) GetAllConsumers(ctx context.Context, pageToken *utils.PaginationToken) ([]*domain.Consumer, error) {
	filter, opts, err := r.buildFilterAndOptions(pageToken.Offset, pageToken.Limit, pageToken.Reversed)
	if err != nil {
		return nil, err
	}

	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	consumers, err := r.decodeConsumers(ctx, cur, pageToken.Reversed)
	if err != nil {
		return nil, err
	}

	return consumers, nil
}

// InsertNewConsumer inserts a new consumer into the MongoDB collection.
func (r *consumerRepository) InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error) {
	result, err := r.collection.InsertOne(ctx, consumer)
	if err != nil {
		return nil, err
	}

	consumer.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return consumer, nil
}

// GetTotalCount return the total collection count using the $collStats
func (r *consumerRepository) GetTotalCount(ctx context.Context) (int, error) {
	count, err := GetTotalCount(ctx, r.collection)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// buildFilterAndOptions constructs the filter and options for the MongoDB query.
func (r *consumerRepository) buildFilterAndOptions(offset string, limit int, reversed bool) (any, *options.FindOptions, error) {
	var filter any = bson.D{}
	opts := options.Find().SetLimit(int64(limit + 1))

	if offset != "" {
		id, err := primitive.ObjectIDFromHex(offset)
		if err != nil {
			return nil, nil, err
		}
		filter = bson.M{"_id": bson.M{"$gt": id}}
	}

	if reversed {
		opts = opts.SetSort(bson.M{"_id": -1})
		if offset != "" {
			id, _ := primitive.ObjectIDFromHex(offset)
			filter = bson.M{"_id": bson.M{"$lt": id}}
		}
	}

	return filter, opts, nil
}

// decodeConsumers decodes the MongoDB cursor into a slice of domain.Consumer.
func (r *consumerRepository) decodeConsumers(ctx context.Context, cur *mongo.Cursor, reversed bool) ([]*domain.Consumer, error) {
	consumers := []*domain.Consumer{}

	for cur.Next(ctx) {
		consumer := &domain.Consumer{}
		err := cur.Decode(consumer)
		if err != nil {
			return nil, err
		}

		if reversed {
			consumers = append([]*domain.Consumer{consumer}, consumers...)
		} else {
			consumers = append(consumers, consumer)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return consumers, nil
}
