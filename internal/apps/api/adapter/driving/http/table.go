package http

import (
	"context"
	"encoding/json"
	"net/http"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/port"
	"reprocess-gui/internal/logger"
)

type tableHandler struct {
	config *config.Config
	log    *logger.Logger
	svc    port.TableService
}

// NewTableHandler creates a new instance of tableHandler.
func NewTableHandler(config *config.Config, log *logger.Logger, svc port.TableService) *tableHandler {
	return &tableHandler{
		config: config,
		log:    log,
		svc:    svc,
	}
}

// GetAllTables retrieves all tables and writes them as JSON to the response.
func (h *tableHandler) GetAllTables(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tables, err := h.svc.GetAllTables(context.Background())
	if err != nil {
		handleError(h.log, w, err, "failed getting tables")
		return
	}

	res, err := json.Marshal(tables)
	if err != nil {
		handleError(h.log, w, err, "failed marshaling tables")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		handleError(h.log, w, err, "failed writing get all tables response to client")
	}
}

// GetTableByTeam retrieves a table by team and writes it as JSON to the response.
func (h *tableHandler) GetTableByTeam(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")

	w.Header().Set("Content-Type", "application/json")

	table, err := h.svc.GetTableByTeam(context.Background(), team)
	if err != nil {
		handleError(h.log, w, err, "failed getting table by team")
		return
	}

	res, err := json.Marshal(table)
	if err != nil {
		handleError(h.log, w, err, "failed marshaling table by team")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		handleError(h.log, w, err, "failed writing get table by team response to client")
	}
}
