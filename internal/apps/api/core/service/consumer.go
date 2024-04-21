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

// GetPagedConsumers retrieves all consumers from the repository.
func (s *consumerService) GetPagedConsumers(ctx context.Context, token string, limit int) (*domain.PagedConsumer, error) {
	pageToken, pagination, err := s.preparePagination(token, limit)
	if err != nil {
		return nil, err
	}

	consumers, err := s.repo.GetPagedConsumers(ctx, pageToken.Offset, pageToken.Limit+1, pageToken.Reversed)
	if err != nil {
		return nil, err
	}

	pagination.TotalRecords, err = s.repo.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}

	if token != "" {
		pagination.PrevPage, consumers, err = utils.GeneratePrevPageToken(*pageToken, consumers, s.config.JWT.Secret)
		if err != nil {
			return nil, err
		}
	}

	pagination.NextPage, consumers, err = utils.GenerateNextPageToken(*pageToken, consumers, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	return &domain.PagedConsumer{
		Consumers:  consumers,
		Pagination: pagination,
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

func (s *consumerService) preparePagination(token string, limit int) (*utils.PaginationToken, *utils.Pagination, error) {
	pageToken := &utils.PaginationToken{}
	pagination := &utils.Pagination{}

	if token != "" {
		err := utils.ParsePaginationToken(token, s.config.JWT.Secret, pageToken)
		if err != nil {
			return nil, nil, err
		}
	}

	if limit > 0 {
		pageToken.Limit = limit
	}

	return pageToken, pagination, nil
}
