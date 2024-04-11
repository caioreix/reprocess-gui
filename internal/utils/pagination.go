package utils

// Pagination represents pagination information.
type Pagination struct {
	TotalRecords int    `json:"total_records,omitempty"`
	NextPage     string `json:"next_page"`
}

// PaginationToken represents a pagination token.
type PaginationToken struct {
	Offset string `json:"offset"`
	Limit  int    `json:"limit"`
}

// GeneratePaginationToken generates a JWT pagination token for the provided object and key.
func GeneratePaginationToken[T any](obj T, key string) (string, error) {
	token, err := NewJWT(obj, key)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ParsePaginationToken parses the JWT pagination token using the provided key into the output object.
func ParsePaginationToken[T any](token, key string, output T) error {
	err := ParseJWT(token, key, output)
	if err != nil {
		return err
	}

	return nil
}
