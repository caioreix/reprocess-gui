package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reprocess-gui/internal/errors"
)

func errorParse(err error) (int, []byte) {
	err = errors.Parse(err)

	switch {
	case errors.Is(err, errors.ErrEmptyResponse):
		return errResponse(http.StatusNoContent, errors.ErrEmptyResponse.Error())
	}

	return errResponse(http.StatusInternalServerError, errors.ErrInternalServerError.Error())
}

func errResponse(status int, message string) (int, []byte) {
	res := struct {
		Error string
	}{message}

	body, err := json.Marshal(res)
	if err != nil {
		return http.StatusInternalServerError, []byte(fmt.Sprintf(`{error: "%s"}`, errors.ErrInternalServerError.Error()))
	}

	return status, body
}

func handleError(w http.ResponseWriter, err error) {
	status, res := errorParse(err)
	w.WriteHeader(status)
	w.Write(res)
}
