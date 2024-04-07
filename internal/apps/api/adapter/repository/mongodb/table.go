package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/logger"
)

type tableRepository struct {
	config     *config.Config
	log        *logger.Logger
	collection *mongo.Collection
}

// NewTableRepository creates a new instance of tableRepository.
func NewTableRepository(config *config.Config, log *logger.Logger, collection *mongo.Collection) *tableRepository {
	return &tableRepository{
		config:     config,
		log:        log,
		collection: collection,
	}
}

// GetAllTables retrieves all tables from the MongoDB collection.
func (r *tableRepository) GetAllTables(ctx context.Context) ([]*domain.Table, error) {
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	tables := []*domain.Table{}
	for cur.Next(ctx) {
		table := &domain.Table{}
		err := cur.Decode(table)
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// GetTableByTeam retrieves a table based on the team name from the MongoDB collection.
func (r *tableRepository) GetTableByTeam(ctx context.Context, team string) (*domain.Table, error) {
	res := r.collection.FindOne(ctx, bson.M{"team": team})
	if err := res.Err(); err != nil {
		return nil, err
	}

	table := &domain.Table{}
	err := res.Decode(table)
	if err != nil {
		return nil, err
	}

	return table, nil
}
