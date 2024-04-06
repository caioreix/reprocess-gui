package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/logger"
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

// InsertNewConsumer inserts a new consumer into the MongoDB collection.
func (r *consumerRepository) InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error) {
	result, err := r.collection.InsertOne(ctx, consumer)
	if err != nil {
		return nil, err
	}

	consumer.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return consumer, nil
}
