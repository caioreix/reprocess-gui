package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/logger"
)

func errorParse(err error, message string) (int, []byte) {
	switch {
	case errors.Is(err, errors.ErrEmptyResponse):
		return errResponse(http.StatusNotFound, fmt.Sprintf("%s: %s", errors.ErrEmptyResponse.Error(), message))
	case errors.Is(err, errors.ErrBadRequest):
		return errResponse(http.StatusBadRequest, fmt.Sprintf("%s: %s", errors.ErrBadRequest.Error(), message))
	}

	return errResponse(http.StatusInternalServerError, errors.ErrInternalServerError.Error())
}

func errResponse(status int, message string) (int, []byte) {
	res := struct {
		Error string `json:"error"`
	}{message}

	body, err := json.Marshal(res)
	if err != nil {
		return http.StatusInternalServerError, []byte(fmt.Sprintf(`{error: "%s"}`, errors.ErrInternalServerError.Error()))
	}

	return status, body
}

func handleError(log *logger.Logger, w http.ResponseWriter, err error, message string) {
	err = errors.Parse(err)
	log.Skip(1).Error(err, message)
	status, res := errorParse(err, message)
	w.WriteHeader(status)
	_, err = w.Write(res)
	if err != nil {
		log.Error(err, "failed writing error response to client")
	}
}
