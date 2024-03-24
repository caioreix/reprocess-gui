package port

import (
	"context"

	"reprocess-gui/internal/apps/api/core/domain"
)

type TableRepository interface {
	GetAllTables(ctx context.Context) ([]*domain.Table, error)
	GetTableByTeam(ctx context.Context, team string) (*domain.Table, error)
}

type TableService interface {
	GetAllTables(ctx context.Context) ([]*domain.Table, error)
	GetTableByTeam(ctx context.Context, team string) (*domain.Table, error)
}
