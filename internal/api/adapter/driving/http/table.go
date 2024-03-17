package http

import (
	"context"
	"encoding/json"
	"net/http"

	"reprocess-gui/internal/api/core/port"
)

type tableHandler struct {
	svc port.TableService
}

func NewTableHandler(svc port.TableService) *tableHandler {
	return &tableHandler{
		svc,
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
