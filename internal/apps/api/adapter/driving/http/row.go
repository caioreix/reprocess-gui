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

type rowHandler struct {
	config *config.Config
	log    *logger.Logger
	svc    port.RowService
}

func NewRowHandler(config *config.Config, log *logger.Logger, svc port.RowService) *rowHandler {
	return &rowHandler{
		config: config,
		log:    log,
		svc:    svc,
	}
}

func (h *rowHandler) InsertNewError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	row := &domain.Row{}
	err := json.NewDecoder(r.Body).Decode(row)
	if err != nil {
		handleError(h.log, w, errors.New(err, errors.ErrBadRequest), "failed decoding row")
		return
	}

	row, err = h.svc.InsertNewError(context.Background(), row)
	if err != nil {
		handleError(h.log, w, err, "failed getting row")
		return
	}

	res, err := json.Marshal(row)
	if err != nil {
		handleError(h.log, w, err, "failed marshaling row")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		handleError(h.log, w, err, "failed writing get all row response to client")
	}
}
