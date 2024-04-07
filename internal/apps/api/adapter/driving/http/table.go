package http

import (
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
func (h *tableHandler) GetAllTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tables, err := h.svc.GetAllTables(r.Context())
	if err != nil {
		handleError(h.log, w, err, "failed getting tables")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tables)
	if err != nil {
		handleError(h.log, w, err, "failed writing get all tables response to client")
	}
}

// GetTableByTeam retrieves a table by team and writes it as JSON to the response.
func (h *tableHandler) GetTableByTeam(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")

	w.Header().Set("Content-Type", "application/json")

	table, err := h.svc.GetTableByTeam(r.Context(), team)
	if err != nil {
		handleError(h.log, w, err, "failed getting table by team")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(table)
	if err != nil {
		handleError(h.log, w, err, "failed writing get table by team response to client")
	}
}
