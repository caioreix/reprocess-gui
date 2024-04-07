package domain

import (
	"time"
)

type status string

const (
	// Pending indicates that the row reprocess is pending.
	Pending status = "pending"
	// Queued indicates that the row is queued for processing.
	Queued status = "queued"
	// Reprocessed indicates that the row has been reprocessed.
	Reprocessed status = "reprocessed"
)

// Row represents the row data.
type Row struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Error     string `json:"error" bson:"error"`
	Value     any    `json:"value" bson:"value"`
	Headers   any    `json:"headers" bson:"headers"`
	TableName string `json:"table_name" bson:"table_name"`
	Target    string `json:"target" bson:"target"`
	Type      string `json:"type" bson:"type"`
	Status    status `json:"status" bson:"status"`

	CreatedAT time.Time `json:"created_at" bson:"created_at"`
	UpdatedAT time.Time `json:"updated_at" bson:"updated_at"`
}
