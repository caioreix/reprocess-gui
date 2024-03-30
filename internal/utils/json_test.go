package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"reprocess-gui/internal/utils"
)

func TestLoadJSONToStruct(t *testing.T) {
	type TestData struct {
		Error string
		Nest  struct {
			Active bool
		}
		Price float64
	}

	jsonFile := "testdata/test.json"
	t.Run("Succeed unmarshal", func(t *testing.T) {
		t.Parallel()
		want := &TestData{
			Error: "My error",
			Nest: struct{ Active bool }{
				Active: true,
			},
			Price: 10.99,
		}

		got := &TestData{}
		b, err := utils.LoadJSONToStruct(jsonFile, got)
		assert.NoError(t, err)
		assert.Equal(t, want, got)

		wantB, err := os.ReadFile(jsonFile)
		require.NoError(t, err)
		assert.Equal(t, string(b), string(wantB))
	})

	t.Run("Failed while reading file", func(t *testing.T) {
		t.Parallel()
		var got *TestData
		b, err := utils.LoadJSONToStruct(":)", got)

		assert.Error(t, err)
		assert.Nil(t, b)
	})

	t.Run("Failed while unmarshal", func(t *testing.T) {
		t.Parallel()
		var got *TestData
		b, err := utils.LoadJSONToStruct(jsonFile, got)

		assert.Error(t, err)
		assert.Nil(t, b)
	})

	t.Run("Not equal unmarshal", func(t *testing.T) {
		t.Parallel()
		want := &TestData{
			Error: "My error",
			Nest: struct{ Active bool }{
				Active: false,
			},
			Price: 10.99,
		}

		got := &TestData{}
		b, err := utils.LoadJSONToStruct(jsonFile, got)
		assert.NoError(t, err)
		assert.NotEqual(t, want, got)

		wantB, err := os.ReadFile(jsonFile)
		require.NoError(t, err)
		assert.Equal(t, b, wantB)
	})
}
