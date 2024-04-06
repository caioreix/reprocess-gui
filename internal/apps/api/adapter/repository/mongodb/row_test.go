package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"reprocess-gui/internal/apps/api/adapter/repository/mongodb"
	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/common"
	"reprocess-gui/internal/logger"
)

func TestInsertNewError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Sucess", func(t *mtest.T) {
		var (
			ctx, config, logger, collection = rowSetupTest(t)
			row                             = &domain.Row{Status: domain.Pending}
		)

		t.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: 1},
			{Key: "acknowledged", Value: 1},
		})

		s := mongodb.NewRowRepository(config, logger, collection)

		newRow, err := s.InsertNewError(ctx, row)
		require.NoError(t, err)
		require.NotEmpty(t, newRow.ID)
	})

	mt.Run("Fail", func(t *mtest.T) {
		var (
			ctx, config, logger, collection = rowSetupTest(t)
			row                             = &domain.Row{Status: domain.Pending}
		)

		t.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		s := mongodb.NewRowRepository(config, logger, collection)

		newRow, err := s.InsertNewError(ctx, row)
		require.Error(t, err)
		require.Empty(t, newRow)
	})
}

func rowSetupTest(t *mtest.T) (context.Context, *config.Config, *logger.Logger, *mongo.Collection) {
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
