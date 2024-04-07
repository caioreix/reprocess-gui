package domain

// Table represents the table data.
type Table struct {
	Name          string          `json:"name" bson:"name"`
	Team          string          `json:"team" bson:"team"`
	CustomColumns []*CustomColumn `json:"custom_columns" bson:"custom_columns"`
	Default       bool            `json:"default" bson:"default"`
}

// CustomColumn represents the table custom columns data.
type CustomColumn struct {
	Name string `json:"name" bson:"name"`
	Path string `json:"path" bson:"path"`
}
