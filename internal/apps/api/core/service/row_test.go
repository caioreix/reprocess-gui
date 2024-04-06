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
)

func TestInsertNewError(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = rowSetupTest(t)
			row                           = &domain.Row{}
			want                          = &domain.Row{}
		)

		repoMock.
			On("InsertNewError", ctx, row).
			Return(want, nil).Once()

		s := service.NewRowService(config, logger, repoMock)

		got, err := s.InsertNewError(ctx, row)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Fail", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = rowSetupTest(t)
			row                           = &domain.Row{}
		)

		repoMock.
			On("InsertNewError", ctx, row).
			Return(nil, errors.ErrEmptyResponse).Once()

		s := service.NewRowService(config, logger, repoMock)

		got, err := s.InsertNewError(ctx, row)
		assert.ErrorIs(t, err, errors.ErrEmptyResponse)
		assert.Nil(t, got)
	})
}

func rowSetupTest(t *testing.T) (context.Context, *config.Config, *logger.Logger, *portmock.RowRepository) {
	t.Helper()

	var (
		ctx       = context.TODO()
		config    = &config.Config{}
		repoMock  = portmock.NewRowRepository(t)
		loggerCfg = logger.Config{Level: "debug", Environment: common.EnvironmentTest}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return ctx, config, logger, repoMock
}
