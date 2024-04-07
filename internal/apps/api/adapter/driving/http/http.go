package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
)

func errorParse(err error, message string) (int, domain.Error) {
	switch {
	case errors.Is(err, errors.ErrEmptyResponse):
		return http.StatusNotFound, domain.Error{Error: fmt.Sprintf("%s: %s", errors.ErrEmptyResponse.Error(), message)}
	case errors.Is(err, errors.ErrBadRequest):
		return http.StatusBadRequest, domain.Error{Error: fmt.Sprintf("%s: %s", errors.ErrBadRequest.Error(), message)}
	}

	return http.StatusInternalServerError, domain.Error{Error: errors.ErrInternalServerError.Error()}
}

func handleError(log *logger.Logger, w http.ResponseWriter, err error, message string) {
	err = errors.Parse(err)
	log.Skip(1).Error(err, message)
	status, res := errorParse(err, message)
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(err, "failed writing error response to client")
	}
}
