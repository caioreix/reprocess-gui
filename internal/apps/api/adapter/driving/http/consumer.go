package http

import (
	"context"
	"encoding/json"
	"net/http"

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

func NewConsumerHandler(config *config.Config, log *logger.Logger, svc port.ConsumerService) *consumerHandler {
	return &consumerHandler{
		config: config,
		log:    log,
		svc:    svc,
	}
}

func (h *consumerHandler) InsertNewConsumer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	consumer := &domain.Consumer{}
	err := json.NewDecoder(r.Body).Decode(consumer)
	if err != nil {
		handleError(h.log, w, errors.New(err, errors.ErrBadRequest), "failed decoding consumer request body")
		return
	}

	consumer, err = h.svc.InsertNewConsumer(context.Background(), consumer)
	if err != nil {
		handleError(h.log, w, err, "failed inserting new consumer")
		return
	}

	res, err := json.Marshal(consumer)
	if err != nil {
		handleError(h.log, w, err, "failed marshaling new consumer")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		handleError(h.log, w, err, "failed writing new consumer response to client")
	}
}
