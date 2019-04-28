package model

import "go-common/library/time"

// Build build.
type Build struct {
	ID       int64     `json:"id"`
	AppID    int64     `json:"app_id"`
	Name     string    `json:"name"`
	TagID    int64     `json:"tag_id"`
	Operator string    `json:"operator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// TableName build.
func (Build) TableName() string {
	return "build"
}
