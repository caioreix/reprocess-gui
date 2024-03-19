package errors

import (
	"errors"
	"strings"
)

var (
	ErrEmptyResponse       = new("empty response value")
	ErrNotMapped           = new("not mapped")
	ErrInternalServerError = new("internal server error")
)

type err struct {
	original error
	message  string
}

func new(message string) error {
	return &err{
		original: errors.New(message),
	}
}

func (e *err) Error() string {
	if e == nil {
		return ""
	}

	msg := e.original.Error()
	if e.message != "" {
		msg += ": " + e.message
	}
	return msg
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Parse(oErr error, messages ...string) error {
	if oErr == nil {
		return nil
	}

	e, ok := oErr.(*err)
	if ok {
		for _, message := range messages {
			if e.message != "" {
				e.message += ": "
			}
			e.message += message
		}
		return e
	}

	switch {
	}

	return &err{
		original: ErrNotMapped,
		message:  strings.Join(messages, ": "),
	}
}
