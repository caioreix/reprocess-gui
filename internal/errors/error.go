package errors

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrEmptyResponse       = new("empty response value")
	ErrNotMapped           = new("not mapped")
	ErrInternalServerError = new("internal server error")

	// HTTP
	ErrBadRequest = new("bad request")
	ErrNotFound   = new("not found")
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

func New(err, target error) error {
	return fmt.Errorf("%w: %w", err, target)
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

func Is(errTree, target error) bool {
	e, ok := errTree.(*err)
	if ok {
		errTree = e.original
	}
	if errors.Is(errTree, target) {
		return true
	}
	t, ok := target.(*err)
	if ok {
		target = t.original
	}

	return errors.Is(errTree, target)
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
	case errors.Is(oErr, mongo.ErrNoDocuments):
		return fmt.Errorf("%w: %w", oErr, ErrEmptyResponse)
	}

	return &err{
		original: fmt.Errorf("%w: %w", oErr, ErrNotMapped),
		message:  strings.Join(messages, ": "),
	}
}
