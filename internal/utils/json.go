package utils

import (
	"encoding/json"
	"os"
)

// LoadJSONToStruct reads JSON data from the file specified by the file parameter
// and unmarshals it into the provided output structure.
// It returns the bytes read from the file and any error encountered.
func LoadJSONToStruct[T any](file string, output T) ([]byte, error) {
	b, err := os.ReadFile(file) // #nosec G304
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, output); err != nil {
		return nil, err
	}

	return b, nil
}
