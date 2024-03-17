package domain

type Table struct {
	Name          string          `json:"name" bson:"name"`
	CustomColumns []*CustomColumn `json:"custom_columns" bson:"custom_columns"`
	Default       bool            `json:"default" bson:"default"`
}

type CustomColumn struct {
	Name string `json:"name" bson:"name"`
	Path string `json:"path" bson:"path"`
}
