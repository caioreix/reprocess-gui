package errors

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrEmptyResponse represents an error indicating an empty response value.
	ErrEmptyResponse = newErr("empty response value")
	// ErrNotMapped represents an error indicating a situation where the error is not mapped.
	ErrNotMapped = newErr("not mapped")
	// ErrInternalServerError represents an internal server error.
	ErrInternalServerError = newErr("internal server error")

	// ErrBadRequest represents a bad request error.
	ErrBadRequest = newErr("bad request")
	// ErrNotFound represents a not found error.
	ErrNotFound = newErr("not found")
)

type err struct {
	original error
	message  string
}

func newErr(message string) error {
	return &err{
		original: errors.New(message),
	}
}

// New returns a formatted error that combines two errors.
func New(err, target error) error {
	return fmt.Errorf("%w: %w", err, target)
}

// Error returns the error message.
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

// Is checks if the error is of the same type as the target error.
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

// Parse handles errors and returns a formatted error.
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
