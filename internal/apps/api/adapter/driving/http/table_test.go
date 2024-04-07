package http_test

import (
	"bytes"
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
	"reprocess-gui/internal/common"
	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
)

func TestGetTableByTeam(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			config, logger, serviceMock = tableSetupTest(t)
			want                        = &domain.Table{
				Name: "table1", Team: "team1", Default: true,
			}
		)
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/tables/team1", nil)
		req.SetPathValue("team", "team1")
		w := httptest.NewRecorder()

		serviceMock.
			On("GetTableByTeam", mock.AnythingOfType("context.backgroundCtx"), "team1").
			Return(want, nil).Once()

		handler := handler.NewTableHandler(config, logger, serviceMock)
		handler.GetTableByTeam(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
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
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
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
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
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
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
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
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})
}

func tableSetupTest(t *testing.T) (*config.Config, *logger.Logger, *portmock.TableService) {
	t.Helper()

	var (
		config      = &config.Config{}
		serviceMock = portmock.NewTableService(t)
		loggerCfg   = logger.Config{Level: "debug", Environment: common.EnvironmentTest}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return config, logger, serviceMock
}
