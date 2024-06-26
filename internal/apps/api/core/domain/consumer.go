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
	ErrorPath  string `json:"error_path" bson:"error_path"`
	RowTTL     int    `json:"row_ttl" bson:"row_ttl"`

	CreatedAT time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAT time.Time `json:"updated_at" bson:"updated_at"`
}

// GetOffset return the table offset to make token pagination
func (c *Consumer) GetOffset() string {
	return c.ID
}

// PagedConsumer represents the consumer with token pagination data.
type PagedConsumer struct {
	Consumers  []*Consumer       `json:"consumers"`
	Pagination *utils.Pagination `json:"pagination" bson:"-"`
}
