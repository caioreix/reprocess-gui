package utils

// Paginable represents an entity that can be paginated.
type Paginable interface {
	GetOffset() string
}

// Pagination represents pagination information.
type Pagination struct {
	TotalRecords int    `json:"total_records,omitempty"`
	NextPage     string `json:"next_page"`
	PrevPage     string `json:"prev_page"`
}

// PaginationToken represents a pagination token.
type PaginationToken struct {
	Offset   string `json:"offset"`
	Limit    int    `json:"limit"`
	Reversed bool   `json:"reversed"`
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

// PreparePagination prepares pagination information based on the provided token, secret, and limit.
func PreparePagination(token, secret string, limit int) (*PaginationToken, *Pagination, error) {
	pageToken := &PaginationToken{}
	pagination := &Pagination{}

	if token != "" {
		err := ParsePaginationToken(token, secret, pageToken)
		if err != nil {
			return nil, nil, err
		}
	}

	if limit > 0 {
		pageToken.Limit = limit
	}

	return pageToken, pagination, nil
}

// GeneratePrevPageToken generates a pagination token for the previous page
// based on the provided PaginationToken, list of Paginable items, and secret.
// It adjusts the PaginationToken accordingly and returns the generated token,
// updated list of Paginable items, and any error encountered during token generation.
//
// ATTENTION: In the repository increment the size of the search limit by one to ensure the next/previous logic work.
func GeneratePrevPageToken[P Paginable](pt PaginationToken, items []P, secret string) (string, []P, error) {
	if pt.Reversed {
		if len(items) <= pt.Limit {
			return "", items, nil
		}
		items = items[1:]
	}

	pt.Offset = items[0].GetOffset()
	pt.Reversed = true

	token, err := GeneratePaginationToken(&pt, secret)
	if err != nil {
		return "", nil, err
	}

	return token, items, nil
}

// GenerateNextPageToken generates a pagination token for the next page
// based on the provided PaginationToken, list of Paginable items, and secret.
// It adjusts the PaginationToken accordingly and returns the generated token,
// updated list of Paginable items, and any error encountered during token generation.
//
// ATTENTION: In the repository increment the size of the search limit by one to ensure the next/previous logic work.
func GenerateNextPageToken[P Paginable](pt PaginationToken, items []P, secret string) (string, []P, error) {
	if !pt.Reversed {
		if len(items) <= pt.Limit {
			return "", items, nil
		}
		items = items[:len(items)-1]
	}

	pt.Offset = items[len(items)-1].GetOffset()
	pt.Reversed = false

	token, err := GeneratePaginationToken(&pt, secret)
	if err != nil {
		return "", nil, err
	}

	return token, items, nil
}
