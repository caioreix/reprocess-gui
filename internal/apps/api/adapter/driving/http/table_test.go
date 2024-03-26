package http_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "reprocess-gui/internal/apps/api/adapter/driving/http"
	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	portmock "reprocess-gui/internal/apps/api/core/port/mocks"
	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
)

func TestGetTableByTeam(t *testing.T) {
	t.Run("Sucess", func(t *testing.T) {
		var (
			config, logger, serviceMock = tableSetupTest(t)
			want                        = &domain.Table{
				Name: "table1", Team: "team1", Default: true,
			}
		)
		expected, err := json.Marshal(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/tables/team1", nil)
		w := httptest.NewRecorder()

		serviceMock.
			On("GetTableByTeam", mock.AnythingOfType("context.backgroundCtx"), "").
			Return(want, nil).Once()

		handler := handler.NewTableHandler(config, logger, serviceMock)
		handler.GetTableByTeam(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestGetAllTables(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			config, logger, serviceMock = tableSetupTest(t)
			want                        = []*domain.Table{
				{Name: "table1", Default: true},
				{Name: "table2"},
			}
		)
		expected, err := json.Marshal(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/tables", nil)
		w := httptest.NewRecorder()

		serviceMock.
			On("GetAllTables", mock.AnythingOfType("context.backgroundCtx")).
			Return(want, nil).Once()

		handler := handler.NewTableHandler(config, logger, serviceMock)
		handler.GetAllTables(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Fail", func(t *testing.T) {
		var (
			config, logger, serviceMock = tableSetupTest(t)
			want                        = struct {
				Error string `json:"error"`
			}{
				Error: "empty response value: failed getting tables",
			}
		)
		expected, err := json.Marshal(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/tables", nil)
		w := httptest.NewRecorder()

		serviceMock.
			On("GetAllTables", mock.AnythingOfType("context.backgroundCtx")).
			Return(nil, errors.ErrEmptyResponse).Once()

		handler := handler.NewTableHandler(config, logger, serviceMock)
		handler.GetAllTables(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func tableSetupTest(t *testing.T) (*config.Config, *logger.Logger, *portmock.TableService) {
	t.Helper()

	var (
		config      = &config.Config{Log: config.Log{Level: "debug"}}
		serviceMock = portmock.NewTableService(t)
	)

	logger, err := logger.New(config)
	require.NoError(t, err)

	return config, logger, serviceMock
}
