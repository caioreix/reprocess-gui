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
	"reprocess-gui/internal/utils"
)

func TestGetAllConsumers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(mt *mtest.T) {
		ctx, config, log, collection := consumerSetupTest(mt)
		want := []*domain.Consumer{
			{Name: "consumer1", Type: "kafka"},
			{Name: "consumer2"},
		}
		pageToken := &utils.PaginationToken{}

		wantConsumers := []primitive.D{}
		for i, c := range want {
			identifier := mtest.NextBatch
			if i == 0 {
				identifier = mtest.FirstBatch
			}
			wantConsumers = append(wantConsumers, mtest.CreateCursorResponse(1, "foo.bar", identifier, bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "name", Value: c.Name},
				{Key: "type", Value: c.Type},
			}))
		}
		wantConsumers = append(wantConsumers, mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch))
		mt.AddMockResponses(wantConsumers...)

		repo := mongodb.NewConsumerRepository(config, log, collection)
		consumers, err := repo.GetAllConsumers(ctx, pageToken)

		assert.Nil(mt, err)
		for _, v := range consumers {
			assert.NotEmpty(mt, v.ID)
			v.ID = ""
		}
		assert.Equal(mt, want, consumers)
	})
}

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
		loggerCfg  = logger.Config{Level: "debug", Environment: common.EnvironmentTest}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return ctx, config, logger, collection
}
