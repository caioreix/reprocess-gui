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
	"reprocess-gui/internal/utils"
)

func TestInsertNewError(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			config, logger, serviceMock = rowSetupTest(t)
			want                        = &domain.Row{}
			row                         = &domain.Row{}
		)
		_, err := utils.LoadJSONToStruct("testdata/row.json", want)
		require.NoError(t, err)
		_, err = utils.LoadJSONToStruct("testdata/row.json", row)
		require.NoError(t, err)
		expected := &bytes.Buffer{}
		err = json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		b := &bytes.Buffer{}
		err = json.NewEncoder(b).Encode(row)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/error", b)
		w := httptest.NewRecorder()

		serviceMock.
			On("InsertNewError", mock.AnythingOfType("context.backgroundCtx"), row).
			Return(want, nil).Once()

		handler := handler.NewRowHandler(config, logger, serviceMock)
		handler.InsertNewError(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})

	t.Run("Fail: bad request", func(t *testing.T) {
		var (
			config, logger, serviceMock = rowSetupTest(t)
			want                        = struct {
				Error string `json:"error"`
			}{
				Error: "bad request: failed decoding row",
			}
		)
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()

		handler := handler.NewRowHandler(config, logger, serviceMock)
		handler.InsertNewError(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})

	t.Run("Fail: not found", func(t *testing.T) {
		var (
			config, logger, serviceMock = rowSetupTest(t)
			want                        = struct {
				Error string `json:"error"`
			}{
				Error: "empty response value: failed getting row",
			}
			row = &domain.Row{}
		)
		_, err := utils.LoadJSONToStruct("testdata/row.json", row)
		require.NoError(t, err)
		expected := &bytes.Buffer{}
		err = json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		b := &bytes.Buffer{}
		err = json.NewEncoder(b).Encode(row)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/error", b)
		w := httptest.NewRecorder()

		serviceMock.
			On("InsertNewError", mock.AnythingOfType("context.backgroundCtx"), row).
			Return(nil, errors.ErrEmptyResponse).Once()

		handler := handler.NewRowHandler(config, logger, serviceMock)
		handler.InsertNewError(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})
}

func rowSetupTest(t *testing.T) (*config.Config, *logger.Logger, *portmock.RowService) {
	t.Helper()

	var (
		config      = &config.Config{}
		serviceMock = portmock.NewRowService(t)
		loggerCfg   = logger.Config{Level: "debug", Environment: common.EnvironmentTest}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return config, logger, serviceMock
}
