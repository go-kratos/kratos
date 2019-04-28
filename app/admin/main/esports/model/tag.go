package model

// Tag .
type Tag struct {
	ID     int64  `json:"id" form:"id"`
	Name   string `json:"name" form:"name" validate:"required"`
	Status int    `json:"status" form:"status"`
}

// TableName .
func (t Tag) TableName() string {
	return "es_tags"
}
