package http

import (
	"context"
	"encoding/json"
	"net/http"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/port"
)

type tableHandler struct {
	config *config.Config
	svc    port.TableService
}

func NewTableHandler(config *config.Config, svc port.TableService) *tableHandler {
	return &tableHandler{
		config: config,
		svc:    svc,
	}
}

func (h *tableHandler) GetAllTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.svc.GetAllTables(context.Background())
	if err != nil {
		return
	}

	res, err := json.Marshal(tables)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
