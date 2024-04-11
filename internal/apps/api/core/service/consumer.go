package service

import (
	"context"
	"time"

	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/domain"
	"reprocess-gui/internal/apps/api/core/port"
	"reprocess-gui/internal/logger"
	"reprocess-gui/internal/utils"
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

// GetAllConsumers retrieves all consumers from the repository.
func (s *consumerService) GetAllConsumers(ctx context.Context, pageToken string, limit int) (*domain.PagedConsumer, error) {
	pagination := &utils.PaginationToken{}
	if pageToken != "" {
		err := utils.ParsePaginationToken(pageToken, s.config.JWT.Secret, pagination)
		if err != nil {
			return nil, err
		}
	}
	if limit > 0 {
		pagination.Limit = limit
	}

	consumers, err := s.repo.GetAllConsumers(ctx, pagination)
	if err != nil {
		return nil, err
	}

	if pagination.Limit >= len(consumers) {
		return &domain.PagedConsumer{
			Consumers:  consumers,
			Pagination: &utils.Pagination{},
		}, nil
	}
	consumers = consumers[:len(consumers)-1]

	pagination.Offset = consumers[len(consumers)-1].ID
	token, err := utils.GeneratePaginationToken(pagination, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	return &domain.PagedConsumer{
		Consumers: consumers,
		Pagination: &utils.Pagination{
			NextPage: token,
		},
	}, nil
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
