package port

import (
	"context"

	"reprocess-gui/internal/apps/api/core/domain"
)

// ConsumerRepository provides methods for interacting with consumer data in the repository.
type ConsumerRepository interface {
	InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error)
}

// ConsumerService provides methods for performing operations related to consumers.
type ConsumerService interface {
	InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error)
}
