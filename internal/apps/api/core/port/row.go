package port

import (
	"context"

	"reprocess-gui/internal/apps/api/core/domain"
)

// RowRepository provides methods for interacting with row data in the repository.
type RowRepository interface {
	InsertNewError(ctx context.Context, row *domain.Row) (*domain.Row, error)
}

// RowService provides methods for performing operations related to rows.
type RowService interface {
	InsertNewError(ctx context.Context, row *domain.Row) (*domain.Row, error)
}
