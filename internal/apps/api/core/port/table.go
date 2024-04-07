package port

import (
	"context"

	"reprocess-gui/internal/apps/api/core/domain"
)

// TableRepository provides methods for interacting with table data in the repository.
type TableRepository interface {
	GetAllTables(ctx context.Context) ([]*domain.Table, error)
	GetTableByTeam(ctx context.Context, team string) (*domain.Table, error)
}

// TableService provides methods for performing operations related to tables.
type TableService interface {
	GetAllTables(ctx context.Context) ([]*domain.Table, error)
	GetTableByTeam(ctx context.Context, team string) (*domain.Table, error)
}
