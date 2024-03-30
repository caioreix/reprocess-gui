package utils

import (
	"encoding/json"
	"os"
)

func LoadJSONToStruct[T any](file string, output T) ([]byte, error) {
	var err error
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, output); err != nil {
		return nil, err
	}

	return b, nil
}
