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

func NewTableHandler(config *config.Config, log *logger.Logger, svc port.TableService) *tableHandler {
	return &tableHandler{
		config: config,
		log:    log,
		svc:    svc,
	}
}

func (h *tableHandler) GetAllTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tables, err := h.svc.GetAllTables(context.Background())
	if err != nil {
		h.handleError(w, err, "failed getting tables")
		return
	}

	res, err := json.Marshal(tables)
	if err != nil {
		h.handleError(w, err, "failed marshaling tables")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *tableHandler) GetTableByTeam(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")

	w.Header().Set("Content-Type", "application/json")

	table, err := h.svc.GetTableByTeam(context.Background(), team)
	if err != nil {
		h.handleError(w, err, "failed getting table by team")
		return
	}

	res, err := json.Marshal(table)
	if err != nil {
		h.handleError(w, err, "failed marshaling table by team")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
