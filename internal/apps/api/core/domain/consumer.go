package domain

import (
	"time"

	"reprocess-gui/internal/utils"
)

// Consumer represents the consumer data.
type Consumer struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Name       string `json:"name" bson:"name"`
	Type       string `json:"type" bson:"type"`
	Table      string `json:"table" bson:"table"`
	Team       string `json:"team" bson:"team"`
	Target     string `json:"target" bson:"target"`
	TargetType string `json:"target_type" bson:"target_type"`

	CreatedAT time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAT time.Time `json:"updated_at" bson:"updated_at"`
}

// PagedConsumer represents the consumer with token pagination data.
type PagedConsumer struct {
	Consumers  []*Consumer       `json:"consumers"`
	Pagination *utils.Pagination `json:"pagination" bson:"-"`
}
