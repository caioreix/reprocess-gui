package port

import (
	"context"

	"reprocess-gui/internal/apps/api/core/domain"
)

// ConsumerRepository provides methods for interacting with consumer data in the repository.
type ConsumerRepository interface {
	GetPagedConsumers(ctx context.Context, offset string, limit int, reversed bool) (pagedConsumer []*domain.Consumer, err error)
	InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error)
	GetTotalCount(ctx context.Context) (int, error)
}

// ConsumerService provides methods for performing operations related to consumers.
type ConsumerService interface {
	GetPagedConsumers(ctx context.Context, pageToken string, limit int) (consumers *domain.PagedConsumer, err error)
	InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error)
}
