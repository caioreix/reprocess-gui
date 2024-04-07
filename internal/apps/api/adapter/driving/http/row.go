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

// NewRowHandler creates a new instance of rowHandler.
func NewRowHandler(config *config.Config, log *logger.Logger, svc port.RowService) *rowHandler {
	return &rowHandler{
		config: config,
		log:    log,
		svc:    svc,
	}
}

// InsertNewError handles HTTP POST requests to insert a new row.
// It decodes the JSON request body into a domain.Row, then inserts it into the database.
// If successful, it responds with the inserted row in JSON format.
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

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(row)
	if err != nil {
		handleError(h.log, w, err, "failed writing get all row response to client")
	}
}
