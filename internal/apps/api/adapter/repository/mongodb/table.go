package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
)

type tableRepository struct {
	config     *config.Config
	collection *mongo.Collection
}

func NewTableRepository(config *config.Config, collection *mongo.Collection) *tableRepository {
	return &tableRepository{
		config:     config,
		collection: collection,
	}
}

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
