package service

import (
	"context"
	"time"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/apps/api/core/port"
	"reprocess-gui/internal/logger"
)

type rowService struct {
	config *config.Config
	log    *logger.Logger
	repo   port.RowRepository
}

// NewRowService creates a new instance of rowService.
func NewRowService(config *config.Config, log *logger.Logger, repo port.RowRepository) *rowService {
	return &rowService{
		config: config,
		log:    log,
		repo:   repo,
	}
}

// InsertNewError inserts a new error row into the repository.
func (s *rowService) InsertNewError(ctx context.Context, row *domain.Row) (*domain.Row, error) {
	row.ID = ""
	row.CreatedAT = time.Now()
	row.UpdatedAT = time.Now()

	row, err := s.repo.InsertNewError(ctx, row)
	if err != nil {
		return nil, err
	}

	return row, nil
}
