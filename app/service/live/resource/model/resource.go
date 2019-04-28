package model

import "time"

// Resource def
type Resource struct {
	ID        int64     `json:"id" gorm:"column:id"`
	Platform  string    `json:"platform" form:"platform"`
	Build     int64     `json:"build" form:"build"`
	LimitType int64     `json:"limit" form:"limit_type"`
	StartTime time.Time `json:"start_time" form:"start_time"`
	EndTime   time.Time `json:"end_time" form:"end_time"`
	Type      string    `json:"type" form:"type"`
	Title     string    `json:"title" form:"title"`
	ImageInfo string    `json:"image_info" form:"image_info"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// TableName resource
func (c Resource) TableName() string {
	return "resource"
}
