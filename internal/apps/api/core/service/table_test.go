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

func TestGetTableByTeam(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = tableSetupTest(t)
			want                          = &domain.Table{
				Name:    "table1",
				Team:    "team1",
				Default: true,
			}
		)

		repoMock.
			On("GetTableByTeam", ctx, "team1").
			Return(want, nil).Once()

		s := service.NewTableService(config, logger, repoMock)

		got, err := s.GetTableByTeam(ctx, "team1")
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Success", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = tableSetupTest(t)
		)

		repoMock.
			On("GetTableByTeam", ctx, "team1").
			Return(nil, errors.ErrEmptyResponse).Once()

		s := service.NewTableService(config, logger, repoMock)

		got, err := s.GetTableByTeam(ctx, "team1")
		assert.ErrorIs(t, err, errors.ErrEmptyResponse)
		assert.Nil(t, got)
	})
}

func TestGetAllTables(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = tableSetupTest(t)
			want                          = []*domain.Table{
				{Name: "table1", Default: true},
				{Name: "table2"},
			}
		)

		repoMock.
			On("GetAllTables", ctx).
			Return(want, nil).Once()

		s := service.NewTableService(config, logger, repoMock)

		got, err := s.GetAllTables(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Success", func(t *testing.T) {
		var (
			ctx, config, logger, repoMock = tableSetupTest(t)
		)

		repoMock.
			On("GetAllTables", ctx).
			Return(nil, errors.ErrEmptyResponse).Once()

		s := service.NewTableService(config, logger, repoMock)

		got, err := s.GetAllTables(ctx)
		assert.ErrorIs(t, err, errors.ErrEmptyResponse)
		assert.Nil(t, got)
	})
}

func tableSetupTest(t *testing.T) (context.Context, *config.Config, *logger.Logger, *portmock.TableRepository) {
	t.Helper()

	var (
		ctx      = context.TODO()
		config   = &config.Config{Log: config.Log{Level: "debug"}}
		repoMock = portmock.NewTableRepository(t)
	)

	logger, err := logger.New(config)
	require.NoError(t, err)

	return ctx, config, logger, repoMock
}
