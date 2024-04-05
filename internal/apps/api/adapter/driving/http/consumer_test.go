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
	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
	"reprocess-gui/internal/utils"
)

func TestInsertNewConsumer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			config, logger, serviceMock = consumerSetupTest(t)
			want                        = &domain.Consumer{}
			consumer                    = &domain.Consumer{}
		)
		_, err := utils.LoadJSONToStruct("testdata/consumer.json", want)
		require.NoError(t, err)
		_, err = utils.LoadJSONToStruct("testdata/consumer.json", consumer)
		require.NoError(t, err)
		expected, err := json.Marshal(want)
		require.NoError(t, err)

		b := &bytes.Buffer{}
		err = json.NewEncoder(b).Encode(consumer)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/error", b)
		w := httptest.NewRecorder()

		serviceMock.
			On("InsertNewConsumer", mock.AnythingOfType("context.backgroundCtx"), consumer).
			Return(want, nil).Once()

		handler := handler.NewConsumerHandler(config, logger, serviceMock)
		handler.InsertNewConsumer(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Fail: bad request", func(t *testing.T) {
		var (
			config, logger, serviceMock = consumerSetupTest(t)
			want                        = struct {
				Error string `json:"error"`
			}{
				Error: "bad request: failed decoding consumer request body",
			}
		)
		expected, err := json.Marshal(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()

		handler := handler.NewConsumerHandler(config, logger, serviceMock)
		handler.InsertNewConsumer(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Fail: not found", func(t *testing.T) {
		var (
			config, logger, serviceMock = consumerSetupTest(t)
			want                        = struct {
				Error string `json:"error"`
			}{
				Error: "empty response value: failed inserting new consumer",
			}
			consumer = &domain.Consumer{}
		)
		_, err := utils.LoadJSONToStruct("testdata/consumer.json", consumer)
		require.NoError(t, err)
		expected, err := json.Marshal(want)
		require.NoError(t, err)

		b := &bytes.Buffer{}
		err = json.NewEncoder(b).Encode(consumer)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/error", b)
		w := httptest.NewRecorder()

		serviceMock.
			On("InsertNewConsumer", mock.AnythingOfType("context.backgroundCtx"), consumer).
			Return(nil, errors.ErrEmptyResponse).Once()

		handler := handler.NewConsumerHandler(config, logger, serviceMock)
		handler.InsertNewConsumer(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func consumerSetupTest(t *testing.T) (*config.Config, *logger.Logger, *portmock.ConsumerService) {
	t.Helper()

	var (
		config      = &config.Config{}
		serviceMock = portmock.NewConsumerService(t)
		loggerCfg   = logger.LoggerConfig{Level: "debug", Environment: "test"}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return config, logger, serviceMock
}
