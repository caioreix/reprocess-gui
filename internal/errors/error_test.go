package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewError(t *testing.T) {
	message := "test error"
	err := newErr(message)

	require.Error(t, err)
	require.Equal(t, message, err.Error())
}

func TestNewErr(t *testing.T) {
	e := newErr("test error")
	assert.Equal(t, "test error", e.Error())
}

// func TestNew(t *testing.T) {
// 	err1 := New("error 1")
// 	err2 := New("error 2")
// 	combined := New(err1, err2)
// 	assert.Equal(t, "error 1: error 2", combined.Error())
// }

func TestErrorNil(t *testing.T) {
	var e *err
	assert.Equal(t, "", e.Error())
}

func TestErrorWithMessage(t *testing.T) {
	e := &err{
		original: newErr("original"),
		message:  "additional",
	}
	assert.Equal(t, "original: additional", e.Error())
}

func TestIs(t *testing.T) {
	target := newErr("target")
	wrapped := &err{original: target}

	assert.True(t, Is(wrapped, target))
	assert.True(t, Is(target, wrapped))
	assert.False(t, Is(newErr("other"), target))
}

func TestParse(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		assert.Nil(t, Parse(nil))
	})

	t.Run("err type", func(t *testing.T) {
		e := &err{original: errors.New("original")}
		parsed := Parse(e, "additional")

		assert.Equal(t, "original: additional", parsed.Error())
	})

	t.Run("mongo.ErrNoDocuments", func(t *testing.T) {
		parsed := Parse(mongo.ErrNoDocuments)
		assert.True(t, errors.Is(parsed, ErrEmptyResponse))
	})

	t.Run("unmapped error", func(t *testing.T) {
		oErr := errors.New("unknown")
		parsed := Parse(oErr, "additional")
		assert.True(t, errors.Is(parsed, ErrNotMapped))
		assert.Equal(t, "unknown: not mapped: additional", parsed.Error())
	})
}

// func TestPredefinedErrors(t *testing.T) {
// 	predefinedErrors := []error{
// 		ErrRedisReaderClientNotInitialized,
// 		ErrRedisWriterClientNotInitialized,
// 		ErrCacheNotFound,
// 		ErrStatusCode,
// 		ErrEmptyResponse,
// 		ErrNotMapped,
// 		ErrInternalServerError,
// 		ErrBadRequest,
// 		ErrNotFound,
// 	}

// 	for _, e := range predefinedErrors {
// 		assert.NotNil(t, e)
// 		assert.NotEmpty(t, e.Error())
// 	}
// }
