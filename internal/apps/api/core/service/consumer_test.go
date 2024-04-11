package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	portmock "reprocess-gui/internal/apps/api/core/port/mocks"
	"reprocess-gui/internal/apps/api/core/service"
	"reprocess-gui/internal/common"
	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
	"reprocess-gui/internal/utils"
)

func TestGetAllConsumers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = consumerSetupTest(t)
			consumers                     = []*domain.Consumer{
				{Name: "consumer1", Type: "kafka"},
				{Name: "consumer2"},
			}
			limit           = 5
			parsedPageToken = &utils.PaginationToken{
				Offset: "1234",
				Limit:  limit,
			}
		)
		pageToken, err := utils.GeneratePaginationToken(parsedPageToken, "")
		require.NoError(t, err)

		want := &domain.PagedConsumer{
			Consumers:  consumers,
			Pagination: &utils.Pagination{},
		}

		repoMock.
			On("GetAllConsumers", ctx, parsedPageToken).
			Return(consumers, nil).Once()

		s := service.NewConsumerService(config, logger, repoMock)

		got, err := s.GetAllConsumers(ctx, pageToken, limit)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Success with next", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = consumerSetupTest(t)
			consumers                     = []*domain.Consumer{
				{Name: "consumer1", Type: "kafka"},
				{Name: "consumer2"},
				{Name: "consumer3"},
				{Name: "consumer3"},
			}
			limit           = 2
			parsedPageToken = &utils.PaginationToken{
				Offset: "1234",
				Limit:  limit,
			}
		)
		pageToken, err := utils.GeneratePaginationToken(parsedPageToken, "")
		require.NoError(t, err)

		want := &domain.PagedConsumer{
			Consumers: consumers[:len(consumers)-1],
			Pagination: &utils.Pagination{
				NextPage: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsaW1pdCI6Miwib2Zmc2V0IjoiIn0.sdqbZSXkM50rqPp157jMFDN4biRrv2CyfDDla9n2SoU",
			},
		}

		repoMock.
			On("GetAllConsumers", ctx, parsedPageToken).
			Return(consumers, nil).Once()

		s := service.NewConsumerService(config, logger, repoMock)

		got, err := s.GetAllConsumers(ctx, pageToken, limit)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Fail", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = consumerSetupTest(t)
			pageToken                     = ""
			limit                         = 5
			parsedPageToken               = &utils.PaginationToken{Limit: limit}
		)

		repoMock.
			On("GetAllConsumers", ctx, parsedPageToken).
			Return(nil, errors.ErrEmptyResponse).Once()

		s := service.NewConsumerService(config, logger, repoMock)

		got, err := s.GetAllConsumers(ctx, pageToken, limit)
		assert.ErrorIs(t, err, errors.ErrEmptyResponse)
		assert.Nil(t, got)
	})
}

func TestInsertNewConsumer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = consumerSetupTest(t)
			consumer                      = &domain.Consumer{}
			want                          = &domain.Consumer{}
		)

		repoMock.
			On("InsertNewConsumer", ctx, consumer).
			Return(want, nil).Once()

		s := service.NewConsumerService(config, logger, repoMock)

		got, err := s.InsertNewConsumer(ctx, consumer)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Fail", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = consumerSetupTest(t)
			consumer                      = &domain.Consumer{}
		)

		repoMock.
			On("InsertNewConsumer", ctx, consumer).
			Return(nil, errors.ErrEmptyResponse).Once()

		s := service.NewConsumerService(config, logger, repoMock)

		got, err := s.InsertNewConsumer(ctx, consumer)
		assert.ErrorIs(t, err, errors.ErrEmptyResponse)
		assert.Nil(t, got)
	})
}

func consumerSetupTest(t *testing.T) (context.Context, *config.Config, *logger.Logger, *portmock.ConsumerRepository) {
	t.Helper()

	var (
		ctx       = context.TODO()
		config    = &config.Config{}
		repoMock  = portmock.NewConsumerRepository(t)
		loggerCfg = logger.Config{Level: "debug", Environment: common.EnvironmentTest}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return ctx, config, logger, repoMock
}
