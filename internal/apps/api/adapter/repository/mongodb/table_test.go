package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"reprocess-gui/internal/apps/api/adapter/repository/mongodb"
	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/common"
	"reprocess-gui/internal/logger"
)

func TestGetTableByTeam(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(mt *mtest.T) {
		ctx, config, log, collection := tableSetupTest(mt)
		want := &domain.Table{
			Name: "table1", Team: "team1", Default: true,
		}

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "name", Value: "table1"},
			{Key: "team", Value: "team1"},
			{Key: "default", Value: true},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "name", Value: "table2"},
			{Key: "team", Value: "team1"},
			{Key: "default", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		repo := mongodb.NewTableRepository(config, log, collection)
		tables, err := repo.GetTableByTeam(ctx, "team1")

		assert.Nil(mt, err)
		assert.Equal(mt, want, tables)
	})
}

func TestGetAllTables(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(mt *mtest.T) {
		ctx, config, log, collection := tableSetupTest(mt)
		want := []*domain.Table{
			{Name: "table1", Default: true},
			{Name: "table2"},
		}

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "name", Value: "table1"},
			{Key: "default", Value: true},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "name", Value: "table2"},
			{Key: "default", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		repo := mongodb.NewTableRepository(config, log, collection)
		tables, err := repo.GetAllTables(ctx)

		assert.Nil(mt, err)
		assert.Equal(mt, want, tables)
	})
}

func tableSetupTest(t *mtest.T) (context.Context, *config.Config, *logger.Logger, *mongo.Collection) {
	t.Helper()

	var (
		ctx        = context.TODO()
		config     = &config.Config{}
		collection = t.Coll
		loggerCfg  = logger.Config{Level: "debug", Environment: common.EnvironmentTest}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return ctx, config, logger, collection
}
