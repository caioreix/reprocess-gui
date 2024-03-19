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
		h.log.Error(err, "failed getting tables")
		handleError(w, err)
		return
	}

	res, err := json.Marshal(tables)
	if err != nil {
		h.log.Error(err, "failed marshaling tables")
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
