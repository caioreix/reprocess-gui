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
	var filter any = bson.D{}
	opts := options.Find().SetLimit(int64(pageToken.Limit + 1))
	if pageToken.Offset != "" {
		id, err := primitive.ObjectIDFromHex(pageToken.Offset)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": bson.M{"$gt": id}}
	}

	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	consumers := []*domain.Consumer{}
	for cur.Next(ctx) {
		consumer := &domain.Consumer{}
		err := cur.Decode(consumer)
		if err != nil {
			return nil, err
		}

		consumers = append(consumers, consumer)
	}
	if err := cur.Err(); err != nil {
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
