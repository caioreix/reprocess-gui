package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	message := "test error"
	err := new(message)

	require.Error(t, err)
	require.Equal(t, message, err.Error())
}
