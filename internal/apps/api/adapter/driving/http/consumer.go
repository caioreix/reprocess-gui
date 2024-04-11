package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/apps/api/core/port"
	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
)

type consumerHandler struct {
	config *config.Config
	log    *logger.Logger
	svc    port.ConsumerService
}

// NewConsumerHandler creates a new instance of consumerHandler.
func NewConsumerHandler(config *config.Config, log *logger.Logger, svc port.ConsumerService) *consumerHandler {
	return &consumerHandler{
		config: config,
		log:    log,
		svc:    svc,
	}
}

// GetAllConsumers retrieves all consumers and writes them as JSON to the response.
func (h *consumerHandler) GetAllConsumers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		pageToken    = r.URL.Query().Get("page_token")
		limitStr     = r.URL.Query().Get("limit")
		defaultLimit = 25
		limit        int
		err          error
	)

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			handleError(h.log, w, errors.New(err, errors.ErrBadRequest), "failed converting limit")
			return
		}
	}
	if limit < 1 {
		limit = defaultLimit
	}

	consumers, err := h.svc.GetAllConsumers(r.Context(), pageToken, limit)
	if err != nil {
		handleError(h.log, w, err, "failed getting all consumers")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(consumers)
	if err != nil {
		handleError(h.log, w, err, "failed writing get all consumers response to client")
	}
}

// InsertNewConsumer handles the HTTP POST request to insert a new consumer.
func (h *consumerHandler) InsertNewConsumer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	consumer := &domain.Consumer{}
	err := json.NewDecoder(r.Body).Decode(consumer)
	if err != nil {
		handleError(h.log, w, errors.New(err, errors.ErrBadRequest), "failed decoding consumer request body")
		return
	}

	consumer, err = h.svc.InsertNewConsumer(r.Context(), consumer)
	if err != nil {
		handleError(h.log, w, err, "failed inserting new consumer")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(consumer)
	if err != nil {
		handleError(h.log, w, err, "failed writing new consumer response to client")
	}
}
