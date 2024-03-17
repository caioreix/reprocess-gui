package service

import (
	"context"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/apps/api/core/port"
)

type tableService struct {
	config *config.Config
	repo   port.TableRepository
}

func NewTableService(config *config.Config, repo port.TableRepository) *tableService {
	return &tableService{
		config: config,
		repo:   repo,
	}
}

func (s *tableService) GetAllTables(ctx context.Context) ([]*domain.Table, error) {
	tables, err := s.repo.GetAllTables(ctx)
	if err != nil {
		return nil, err
	}

	return tables, nil
}
