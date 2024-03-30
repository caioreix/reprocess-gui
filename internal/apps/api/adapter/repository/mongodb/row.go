package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/logger"
)

type rowRepository struct {
	config     *config.Config
	log        *logger.Logger
	collection *mongo.Collection
}

func NewRowRepository(config *config.Config, log *logger.Logger, collection *mongo.Collection) *rowRepository {
	return &rowRepository{
		config:     config,
		log:        log,
		collection: collection,
	}
}

func (r *rowRepository) InsertNewError(ctx context.Context, row *domain.Row) (*domain.Row, error) {
	result, err := r.collection.InsertOne(ctx, row)
	if err != nil {
		return nil, err
	}

	row.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return row, nil
}
