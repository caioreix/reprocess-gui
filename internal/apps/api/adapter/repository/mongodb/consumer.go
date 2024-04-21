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

// GetPagedConsumers retrieves all consumers from the MongoDB collection.
func (r *consumerRepository) GetPagedConsumers(ctx context.Context, offset string, limit int, reversed bool) ([]*domain.Consumer, error) {
	filter, opts, err := buildPaginationFilterAndOptions("_id", offset, limit, reversed)
	if err != nil {
		return nil, err
	}

	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	consumers, err := decodePaginationResponse[domain.Consumer](ctx, cur, reversed)
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
