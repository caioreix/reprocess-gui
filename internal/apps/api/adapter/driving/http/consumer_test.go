package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestGetAllConsumers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			config, logger, serviceMock = consumerSetupTest(t)
			want                        = &domain.PagedConsumer{
				Consumers: []*domain.Consumer{
					{Name: "consumer1", Type: "kafka"},
					{Name: "consumer2"},
				},
				Pagination: &utils.Pagination{},
			}
			pageToken = "1234"
			limit     = 25
		)
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/consumers", nil)
		values := req.URL.Query()
		values.Add("page_token", pageToken)
		values.Add("limit", "")
		req.URL.RawQuery = values.Encode()

		w := httptest.NewRecorder()

		serviceMock.
			On("GetAllConsumers", mock.AnythingOfType("context.backgroundCtx"), pageToken, limit).
			Return(want, nil).Once()

		handler := handler.NewConsumerHandler(config, logger, serviceMock)
		handler.GetAllConsumers(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})

	t.Run("Fail invalid limit", func(t *testing.T) {
		var (
			config, logger, serviceMock = consumerSetupTest(t)
			want                        = struct {
				Error string `json:"error"`
			}{
				Error: "bad request: failed converting limit",
			}
			pageToken = "1234"
		)
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/consumers", nil)
		values := req.URL.Query()
		values.Add("page_token", pageToken)
		values.Add("limit", "x")
		req.URL.RawQuery = values.Encode()

		w := httptest.NewRecorder()

		handler := handler.NewConsumerHandler(config, logger, serviceMock)
		handler.GetAllConsumers(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})

	t.Run("Fail empty response", func(t *testing.T) {
		var (
			config, logger, serviceMock = consumerSetupTest(t)
			want                        = struct {
				Error string `json:"error"`
			}{
				Error: "empty response value: failed getting all consumers",
			}
			pageToken = "1234"
			limit     = 10
		)
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/consumers", nil)
		values := req.URL.Query()
		values.Add("page_token", pageToken)
		values.Add("limit", fmt.Sprint(limit))
		req.URL.RawQuery = values.Encode()

		w := httptest.NewRecorder()

		serviceMock.
			On("GetAllConsumers", mock.AnythingOfType("context.backgroundCtx"), pageToken, limit).
			Return(nil, errors.ErrEmptyResponse).Once()

		handler := handler.NewConsumerHandler(config, logger, serviceMock)
		handler.GetAllConsumers(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})
}

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
		expected := &bytes.Buffer{}
		err = json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		b := &bytes.Buffer{}
		err = json.NewEncoder(b).Encode(consumer)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/consumer", b)
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
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
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
		expected := &bytes.Buffer{}
		err := json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/consumer", nil)
		w := httptest.NewRecorder()

		handler := handler.NewConsumerHandler(config, logger, serviceMock)
		handler.InsertNewConsumer(w, req)

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
		expected := &bytes.Buffer{}
		err = json.NewEncoder(expected).Encode(want)
		require.NoError(t, err)

		b := &bytes.Buffer{}
		err = json.NewEncoder(b).Encode(consumer)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/consumer", b)
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
		assert.Equal(t, expected.String(), string(data))
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})
}

func consumerSetupTest(t *testing.T) (*config.Config, *logger.Logger, *portmock.ConsumerService) {
	t.Helper()

	var (
		config      = &config.Config{}
		serviceMock = portmock.NewConsumerService(t)
		loggerCfg   = logger.Config{Level: "debug", Environment: common.EnvironmentTest}
	)

	logger, err := loggerCfg.New()
	require.NoError(t, err)

	return config, logger, serviceMock
}
