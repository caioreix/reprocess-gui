package port

import (
	"context"

	"reprocess-gui/internal/apps/api/core/domain"
)

type RowRepository interface {
	InsertNewError(ctx context.Context, row *domain.Row) (*domain.Row, error)
}

type RowService interface {
	InsertNewError(ctx context.Context, row *domain.Row) (*domain.Row, error)
}
