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
	"reprocess-gui/internal/logger"
)

func TestInsertNewConsumer(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Sucess", func(t *mtest.T) {
		var (
			ctx, config, logger, collection = consumerSetupTest(t)
			consumer                        = &domain.Consumer{Team: "channels"}
		)

		t.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: 1},
			{Key: "acknowledged", Value: 1},
		})

		s := mongodb.NewConsumerRepository(config, logger, collection)

		newConsumer, err := s.InsertNewConsumer(ctx, consumer)
		require.NoError(t, err)
		require.NotEmpty(t, newConsumer.ID)
	})

	mt.Run("Fail", func(t *mtest.T) {
		var (
			ctx, config, logger, collection = consumerSetupTest(t)
			consumer                        = &domain.Consumer{Team: "channels"}
		)

		t.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		s := mongodb.NewConsumerRepository(config, logger, collection)

		newConsumer, err := s.InsertNewConsumer(ctx, consumer)
		require.Error(t, err)
		require.Empty(t, newConsumer)
	})
}

func consumerSetupTest(t *mtest.T) (context.Context, *config.Config, *logger.Logger, *mongo.Collection) {
	t.Helper()

	var (
		ctx        = context.TODO()
		config     = &config.Config{}
		collection = t.Coll
		loggerCfg  = logger.LoggerConfig{Level: "debug", Environment: "test"}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return ctx, config, logger, collection
}
