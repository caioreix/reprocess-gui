package service

import (
	"context"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/apps/api/core/port"
	"reprocess-gui/internal/logger"
)

type tableService struct {
	config *config.Config
	log    *logger.Logger
	repo   port.TableRepository
}

func NewTableService(config *config.Config, log *logger.Logger, repo port.TableRepository) *tableService {
	return &tableService{
		config: config,
		log:    log,
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
