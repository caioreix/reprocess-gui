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
	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
)

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
		ctx      = context.TODO()
		config   = &config.Config{}
		repoMock = portmock.NewConsumerRepository(t)
	)

	logger, err := logger.New("debug")
	require.NoError(t, err)

	return ctx, config, logger, repoMock
}
