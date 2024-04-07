package service

import (
	"context"
	"time"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/apps/api/core/port"
	"reprocess-gui/internal/logger"
)

type consumerService struct {
	config *config.Config
	log    *logger.Logger
	repo   port.ConsumerRepository
}

// NewConsumerService creates a new instance of consumerService.
func NewConsumerService(config *config.Config, log *logger.Logger, repo port.ConsumerRepository) *consumerService {
	return &consumerService{
		config: config,
		log:    log,
		repo:   repo,
	}
}

// InsertNewConsumer inserts a new consumer into the repository.
func (s *consumerService) InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error) {
	consumer.ID = ""
	consumer.CreatedAT = time.Now()
	consumer.UpdatedAT = time.Now()

	consumer, err := s.repo.InsertNewConsumer(ctx, consumer)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
