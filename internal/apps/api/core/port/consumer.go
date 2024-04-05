package port

import (
	"context"

	"reprocess-gui/internal/apps/api/core/domain"
)

type ConsumerRepository interface {
	InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error)
}

type ConsumerService interface {
	InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error)
}
