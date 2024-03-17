package service

import (
	"context"

	"reprocess-gui/internal/api/core/domain"
	"reprocess-gui/internal/api/core/port"
)

type tableService struct {
	repo port.TableRepository
}

func NewTableService(repo port.TableRepository) *tableService {
	return &tableService{
		repo: repo,
	}
}

func (s *tableService) GetAllTables(ctx context.Context) ([]*domain.Table, error) {
	tables, err := s.repo.GetAllTables(ctx)
	if err != nil {
		return nil, err
	}

	return tables, nil
}
